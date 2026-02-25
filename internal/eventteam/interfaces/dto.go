package interfaces

import "time"

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required,min=2,max=40"`
}

type JoinTeamRequest struct {
	TeamName   string `json:"team_name"   validate:"required"`
	JoinSecret string `json:"join_secret" validate:"required"`
}

type EventTeamResponse struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	MemberCount int    `json:"member_count,omitempty"`
}

type CreateTeamResponse struct {
	EventTeamResponse
	JoinSecret string `json:"join_secret"`
}

type EventTeamMemberResponse struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
