package models

type Player struct {
	PlayerID          int    `json:"player_id"`
	Name              string `json:"name"`
	Position          string `json:"position"`
	Overall           int    `json:"overall"`
	GamesPlayed       int    `json:"games_played"`
	GamesStarted      int    `json:"games_started"`
	SnapsPlayed       int    `json:"snaps_played"`
	ClassYear         string `json:"class_year"`       // Freshman, Sophomore, etc.
	RecruitmentYear   int    `json:"recruitment_year"` // Ano de recrutamento
	TeamID            int    `json:"team_id"`
	RecruitmentSource string `json:"recruitment_source"` // Fonte de recrutamento
}

type PlayerGameStats struct {
	PlayerID      int `json:"player_id"`
	ScheduleID    int `json:"schedule_id"`
	Completions   int `json:"completions"`
	PassAttempts  int `json:"pass_attempts"`
	PassingYards  int `json:"passing_yards"`
	PassingTDs    int `json:"passing_tds"`
	Interceptions int `json:"interceptions"`
	RushAttempts  int `json:"rush_attempts"`
	RushingYards  int `json:"rushing_yards"`
	RushingTDs    int `json:"rushing_tds"`
}
