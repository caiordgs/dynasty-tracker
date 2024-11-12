package models

type TeamAssignment struct {
	TeamID  int    `json:"team_id"`
	CoachID int    `json:"coach_id"`
	Year    int    `json:"year"`
	Role    string `json:"role"` // Ex: HC (Head Coach), OC (Offensive Coordinator), etc.
}
