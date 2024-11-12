package models

type HistoricalRecord struct {
	RecordID             int      `json:"record_id"`
	School               string   `json:"school"`
	PlayerName           string   `json:"player_name"`
	YearStart            int      `json:"year_start"`
	YearEnd              int      `json:"year_end"`
	Completions          *int     `json:"completions"`
	Attempts             *int     `json:"attempts"`
	CompletionPercentage *float64 `json:"completion_percentage"`
	PassingYards         *int     `json:"passing_yards"`
	YardsPerAttempt      *float64 `json:"yards_per_attempt"`
	Touchdowns           *int     `json:"touchdowns"`
	Interceptions        *int     `json:"interceptions"`
	PasserRating         *float64 `json:"passer_rating"`
	RushAttempts         *int     `json:"rush_attempts"`
	RushYards            *int     `json:"rush_yards"`
	YardsPerCarry        *float64 `json:"yards_per_carry"`
	RushTDs              *int     `json:"rush_tds"`
	Receptions           *int     `json:"receptions"`
	ReceivingYards       *int     `json:"receiving_yards"`
	YardsPerCatch        *float64 `json:"yards_per_catch"`
	ReceivingTDs         *int     `json:"receiving_tds"`
	PlaysFromScrimmage   *int     `json:"plays_from_scrimmage"`
	YardsFromScrimmage   *int     `json:"yards_from_scrimmage"`
	AvgYardsPerPlay      *float64 `json:"avg_yards_per_play"`
	ScrimmageTDs         *int     `json:"scrimmage_tds"`
}
