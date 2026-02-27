package domain

type InvitePayload struct {
	TeamID string `json:"team_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
