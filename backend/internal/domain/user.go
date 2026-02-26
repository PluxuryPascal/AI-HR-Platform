package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	TeamID       string    `json:"team_id"`
	TeamName     string    `json:"team_name"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Locale       string    `json:"locale"`
}

type RegisterOwnerRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	TeamName  string `json:"team_name"`
}

type CreateUserParams struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	TeamID    string `json:"team_id"`
	Role      string `json:"role"`
}

type Session struct {
	UserID string `json:"user_id"`
	TeamID string `json:"team_id"`
	Role   string `json:"role"`
}
