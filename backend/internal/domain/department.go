package domain

// Department represents the hiring.t_departments table
type Department struct {
	ID     string `json:"id" db:"id"`
	TeamID string `json:"team_id" db:"team_id"`
	Name   string `json:"name" db:"name"`
}

type CreateDepartmentParams struct {
	TeamID string `json:"team_id"`
	Name   string `json:"name"`
}
