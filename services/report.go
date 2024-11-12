package services

import (
	"dynastyTracker/database"
	"fmt"
)

// Estrutura para representar o relatório de desempenho do time
type TeamPerformanceReport struct {
	Year   int `json:"year"`
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
}

// Função que calcula o número de vitórias e derrotas por ano
func GetTeamPerformance() ([]TeamPerformanceReport, error) {
	query := `
        SELECT year,
            SUM(CASE WHEN result = 'Win' THEN 1 ELSE 0 END) AS wins,
            SUM(CASE WHEN result = 'Loss' THEN 1 ELSE 0 END) AS losses
        FROM schedule
        GROUP BY year
        ORDER BY year;
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reports []TeamPerformanceReport
	for rows.Next() {
		var report TeamPerformanceReport
		err := rows.Scan(&report.Year, &report.Wins, &report.Losses)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// Estrutura para representar o relatório de estatísticas do jogador
type PlayerStatsReport struct {
	PlayerName          string `json:"player_name"`
	Position            string `json:"position"`
	PassingYards        int    `json:"passing_yards,omitempty"`
	PassingTouchdowns   int    `json:"passing_touchdowns,omitempty"`
	Interceptions       int    `json:"interceptions,omitempty"`
	RushingYards        int    `json:"rushing_yards,omitempty"`
	RushingTouchdowns   int    `json:"rushing_touchdowns,omitempty"`
	ReceivingYards      int    `json:"receiving_yards,omitempty"`
	ReceivingTouchdowns int    `json:"receiving_touchdowns,omitempty"`
}

// Função que retorna estatísticas dos jogadores para uma posição específica
func GetPlayerStatsByPosition(position string) ([]PlayerStatsReport, error) {
	query := `
        SELECT name, position,
            SUM(passing_yards) AS passing_yards,
            SUM(passing_touchdowns) AS passing_touchdowns,
            SUM(interceptions) AS interceptions,
            SUM(rushing_yards) AS rushing_yards,
            SUM(rushing_touchdowns) AS rushing_touchdowns,
            SUM(receiving_yards) AS receiving_yards,
            SUM(receiving_touchdowns) AS receiving_touchdowns
        FROM players
        JOIN player_game_stats ON players.player_id = player_game_stats.player_id
        WHERE position = ?
        GROUP BY name, position
        ORDER BY name;
    `

	rows, err := database.DB.Query(query, position)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reports []PlayerStatsReport
	for rows.Next() {
		var report PlayerStatsReport
		err := rows.Scan(
			&report.PlayerName, &report.Position,
			&report.PassingYards, &report.PassingTouchdowns, &report.Interceptions,
			&report.RushingYards, &report.RushingTouchdowns,
			&report.ReceivingYards, &report.ReceivingTouchdowns,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}
