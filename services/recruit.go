package services

import (
	"dynastyTracker/database"
	"dynastyTracker/models"
	"fmt"
)

// Função para adicionar um recruta à tabela recruits
func AddRecruit(recruit models.Recruit) error {
	query := `
        INSERT INTO recruits (player_name, class, position, tendency, position_rank, national_rank, stars, hometown, home_state, height, weight, dev_trait, overall, gem_bust, recruitment_source, recruitment_year, team_id)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	_, err := database.DB.Exec(query, recruit.PlayerName, recruit.Class, recruit.Position, recruit.Tendency, recruit.PositionRank, recruit.NationalRank, recruit.Stars, recruit.Hometown, recruit.HomeState, recruit.Height, recruit.Weight, recruit.DevTrait, recruit.Overall, recruit.GemBust, recruit.RecruitmentSource, recruit.RecruitmentYear, recruit.TeamID)
	if err != nil {
		fmt.Printf("Erro ao adicionar recruta: %v\n", err)
		return err
	}
	return nil
}
