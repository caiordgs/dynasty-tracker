package models

type Schedule struct {
	ID              int    `json:"id"`
	TeamID          int    `json:"team_id"`
	TeamName        string `json:"team_name"` // Verifique esta linha
	Year            int    `json:"year"`
	Week            int    `json:"week"`
	Opponent        string `json:"opponent"`
	TeamRanking     int    `json:"team_ranking"`
	OpponentRanking int    `json:"opponent_ranking"`
	TeamPoints      int    `json:"team_points"`
	OpponentPoints  int    `json:"opponent_points"`
	Result          string `json:"result"`
	Site            string `json:"site"`
}
