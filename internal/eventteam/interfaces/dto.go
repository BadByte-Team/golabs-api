package interfaces

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type JoinTeamRequest struct {
	TeamName   string `json:"team_name"`
	JoinSecret string `json:"join_secret"`
}

type EventTeamResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type CreateTeamResponse struct {
	EventTeamResponse
	JoinSecret string `json:"join_secret"`
}
