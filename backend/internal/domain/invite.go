package domain

import "time"

// Invite represents a row in auth.t_invites.
// It is created when an owner sends an invitation email and destroyed
// (cascading its job-access records) when the invitation is accepted.
type Invite struct {
	ID        string    `db:"id"`
	TeamID    string    `db:"team_id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

// CreateInviteParams is the input DTO for the InviteUser use-case method.
type CreateInviteParams struct {
	Email  string
	Role   string
	JobIDs *[]string
}

// InviteRegisterDTO is the public-facing payload returned by ValidateInvite.
// It contains only the information the frontend needs to pre-fill the
// registration form.
type InviteRegisterDTO struct {
	Email string `json:"email"`
}
