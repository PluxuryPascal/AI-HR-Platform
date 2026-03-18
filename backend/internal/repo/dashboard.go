package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"
	"time"
)

type DashboardRepository interface {
	GetStats(ctx context.Context, teamID string) (*domain.DashboardStats, error)
	GetApplicationDynamics(ctx context.Context, teamID string, startDate, endDate time.Time, groupBy string) ([]domain.ChartDataPoint, error)
	GetRecentActivity(ctx context.Context, teamID string, limit int) ([]domain.ActivityLogEntry, error)
}

type dashboardRepo struct {
	dbClient *db.PostgresClient
}

func NewDashboardRepo(dbClient *db.PostgresClient) DashboardRepository {
	return &dashboardRepo{dbClient: dbClient}
}

func (r *dashboardRepo) GetStats(ctx context.Context, teamID string) (*domain.DashboardStats, error) {
	const query = `
		WITH candidates_count AS (
			SELECT COUNT(c.id) as total
			FROM hiring.t_candidates c
			JOIN hiring.t_jobs j ON c.job_id = j.id
			WHERE j.team_id = $1
		),
		jobs_count AS (
			SELECT COUNT(id) as active_jobs
			FROM hiring.t_jobs
			WHERE team_id = $1 AND status = 'status_published'
		),
		interviews_count AS (
			SELECT COUNT(cs.candidate_id) as upcoming_interviews
			FROM hiring.t_candidate_stages cs
			JOIN hiring.t_pipeline_stages ps ON cs.stage_id = ps.id
			JOIN hiring.t_candidates c ON cs.candidate_id = c.id
			JOIN hiring.t_jobs j ON c.job_id = j.id
			WHERE j.team_id = $1 AND (ps.title ILIKE '%Interview%' OR ps.title ILIKE '%Интервью%')
		),
		avg_time AS (
			SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (cs.moved_at - c.created_at)) / 86400), 0) as avg_days
			FROM hiring.t_candidate_stages cs
			JOIN hiring.t_pipeline_stages ps ON cs.stage_id = ps.id
			JOIN hiring.t_candidates c ON cs.candidate_id = c.id
			JOIN hiring.t_jobs j ON c.job_id = j.id
			WHERE j.team_id = $1 AND ps.is_terminal = true
		)
		SELECT 
			(SELECT total FROM candidates_count) as total_candidates,
			(SELECT active_jobs FROM jobs_count) as active_jobs,
			(SELECT upcoming_interviews FROM interviews_count) as upcoming_interviews,
			(SELECT avg_days FROM avg_time) as avg_time_to_hire_days
	`

	var stats domain.DashboardStats
	err := r.dbClient.Pool.QueryRow(ctx, query, teamID).Scan(
		&stats.TotalCandidates,
		&stats.ActiveJobs,
		&stats.UpcomingInterviews,
		&stats.AvgTimeToHireDays,
	)
	if err != nil {
		return nil, fmt.Errorf("dashboard query stats: %w", err)
	}

	// We can compute deltas separately if needed or just return 0 for now.
	return &stats, nil
}

func (r *dashboardRepo) GetApplicationDynamics(ctx context.Context, teamID string, startDate, endDate time.Time, groupBy string) ([]domain.ChartDataPoint, error) {
	truncPart := "month"
	if groupBy == "daily" {
		truncPart = "day"
	}

	query := fmt.Sprintf(`
		SELECT DATE_TRUNC('%s', c.created_at) as label_date, COUNT(*) as count
		FROM hiring.t_candidates c
		JOIN hiring.t_jobs j ON c.job_id = j.id
		WHERE j.team_id = $1 AND c.created_at >= $2 AND c.created_at <= $3
		GROUP BY label_date
		ORDER BY label_date ASC
	`, truncPart)

	rows, err := r.dbClient.Pool.Query(ctx, query, teamID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("dashboard application dynamics: %w", err)
	}
	defer rows.Close()

	var data []domain.ChartDataPoint
	for rows.Next() {
		var pt domain.ChartDataPoint
		if err := rows.Scan(&pt.Date, &pt.Count); err != nil {
			return nil, err
		}
		data = append(data, pt)
	}
	return data, nil
}

func (r *dashboardRepo) GetRecentActivity(ctx context.Context, teamID string, limit int) ([]domain.ActivityLogEntry, error) {
	query := `
		SELECT 
			l.id,
			l.actor_type,
			l.actor_id,
			a.code as action_code,
			l.target_id,
			l.created_at,
			c.first_name,
			c.last_name,
			j.title as job_title
		FROM logging.t_activity_logs l
		JOIN logging.t_action_types a ON l.action_id = a.id
		LEFT JOIN hiring.t_candidates c ON l.target_id = c.id
		LEFT JOIN hiring.t_jobs j ON l.target_id = j.id OR c.job_id = j.id
		WHERE l.team_id = $1
		ORDER BY l.created_at DESC
		LIMIT $2
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, teamID, limit)
	if err != nil {
		return nil, fmt.Errorf("dashboard recent activity: %w", err)
	}
	defer rows.Close()

	var logs []domain.ActivityLogEntry
	for rows.Next() {
		var l domain.ActivityLogEntry
		if err := rows.Scan(
			&l.LogID, &l.ActorType, &l.ActorID, &l.ActionCode, &l.TargetID, &l.CreatedAt,
			&l.CandidateFirstName, &l.CandidateLastName, &l.JobTitle,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
