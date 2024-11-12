package models

// Definindo o modelo do time com todas as colunas
type Team struct {
	TeamID          int    `json:"team_id"`
	Year            int    `json:"year"`
	Role            string `json:"role"`
	School          string `json:"school"`
	Mascot          string `json:"mascot"`
	Abbreviation    string `json:"abbreviation"`
	AltName1        string `json:"alt_name1"`
	Color           string `json:"color"`
	AltColor        string `json:"alt_color"`
	Logo1           string `json:"logo_1"`
	Logo2           string `json:"logo_2"`
	Twitter         string `json:"twitter"`
	LocationVenueID string `json:"location_venue_id"`
	LocationName    string `json:"location_name"`
	LocationCity    string `json:"location_city"`
	LocationState   string `json:"location_state"`
}
