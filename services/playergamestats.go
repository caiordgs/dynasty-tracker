package services

import (
	"dynastyTracker/database"
	"dynastyTracker/models"
	"fmt"
)

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

// Função para adicionar estatísticas de jogo para um jogador
func AddPlayerGameStats(stats models.PlayerGameStats) error {
	query := `
        INSERT INTO playergamestats (player_id, schedule_id, completions, pass_attempts, passing_yards, passing_tds, interceptions, rush_attempts, rushing_yards, rushing_tds)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := database.DB.Exec(query, stats.PlayerID, stats.ScheduleID, stats.Completions, stats.PassAttempts, stats.PassingYards, stats.PassingTDs, stats.Interceptions, stats.RushAttempts, stats.RushingYards, stats.RushingTDs)
	if err != nil {
		fmt.Printf("Erro ao adicionar estatísticas do jogo: %v\n", err)
		return err
	}
	return nil
}
