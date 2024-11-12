package models

type Recruit struct {
	RecruitID         int    `json:"recruit_id"`
	PlayerName        string `json:"player_name"`
	Class             string `json:"class"`
	Position          string `json:"position"`
	Tendency          string `json:"tendency"`
	PositionRank      int    `json:"position_rank"`
	NationalRank      int    `json:"national_rank"`
	Stars             int    `json:"stars"`
	Hometown          string `json:"hometown"`
	HomeState         string `json:"home_state"`
	Height            int    `json:"height"`
	Weight            int    `json:"weight"`
	DevTrait          string `json:"dev_trait"`
	Overall           int    `json:"overall"`
	GemBust           string `json:"gem_bust"`
	RecruitmentSource string `json:"recruitment_source"` // Fonte de recrutamento: "High School" ou "Transfer Portal"
	RecruitmentYear   int    `json:"recruitment_year"`   // Ano de recrutamento
	TeamID            int    `json:"team_id"`            // ID do time
}
