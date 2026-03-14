package domain

import (
	"time"
)

// HiringForecast represents the ai_engine.t_hiring_forecasts table
type HiringForecast struct {
	ID              string     `json:"id" db:"id"`
	TeamID          string     `json:"team_id" db:"team_id"`
	ForecastMonth   time.Time  `json:"forecast_month" db:"forecast_month"`
	PredictedVolume int        `json:"predicted_volume" db:"predicted_volume"`
	GeneratedAt     *time.Time `json:"generated_at,omitempty" db:"generated_at"`
}
