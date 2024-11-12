package services

import (
	"dynastyTracker/database"
	"fmt"
)

// Estrutura para representar o relatório de desempenho do time
type TeamPerformanceReport struct {
	Year               int `json:"year"`
	Wins               int `json:"wins"`
	Losses             int `json:"losses"`
	TotalPointsScored  int `json:"total_points_scored"`
	TotalPointsAllowed int `json:"total_points_allowed"`
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

func GetTeamPerformanceBySeason() ([]TeamPerformanceReport, error) {
	query := `
        SELECT year,
            SUM(CASE WHEN result = 'Win' THEN 1 ELSE 0 END) AS wins,
            SUM(CASE WHEN result = 'Loss' THEN 1 ELSE 0 END) AS losses,
            SUM(team_points) AS total_points_scored,
            SUM(opponent_points) AS total_points_allowed
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
		err := rows.Scan(&report.Year, &report.Wins, &report.Losses, &report.TotalPointsScored, &report.TotalPointsAllowed)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

type GameSummaryReport struct {
	Week           int    `json:"week"`
	TeamName       string `json:"team_name"`
	Opponent       string `json:"opponent"`
	TeamPoints     int    `json:"team_points"`
	OpponentPoints int    `json:"opponent_points"`
	Result         string `json:"result"`
	Site           string `json:"site"`
}

func GetSeasonSummary(year int) ([]GameSummaryReport, error) {
	query := `
        SELECT week, team_name, opponent, team_points, opponent_points, result, site
        FROM schedule
        WHERE year = ?
        ORDER BY week;
    `

	rows, err := database.DB.Query(query, year)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reports []GameSummaryReport
	for rows.Next() {
		var report GameSummaryReport
		err := rows.Scan(&report.Week, &report.TeamName, &report.Opponent, &report.TeamPoints, &report.OpponentPoints, &report.Result, &report.Site)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

type TeamSeasonStats struct {
	Year               int `json:"year"`
	Wins               int `json:"wins"`
	Losses             int `json:"losses"`
	TotalPointsScored  int `json:"total_points_scored"`
	TotalPointsAllowed int `json:"total_points_allowed"`
}

func GetTeamSeasonComparison(teamID int) ([]TeamSeasonStats, error) {
	query := `
        SELECT 
            year,
            SUM(CASE WHEN result = 'Win' THEN 1 ELSE 0 END) AS wins,
            SUM(CASE WHEN result = 'Loss' THEN 1 ELSE 0 END) AS losses,
            SUM(team_points) AS total_points_scored,
            SUM(opponent_points) AS total_points_allowed
        FROM schedule
        WHERE team_id = ?
        GROUP BY year
        ORDER BY year;
    `

	rows, err := database.DB.Query(query, teamID)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var seasonStats []TeamSeasonStats
	for rows.Next() {
		var stats TeamSeasonStats
		err := rows.Scan(
			&stats.Year,
			&stats.Wins,
			&stats.Losses,
			&stats.TotalPointsScored,
			&stats.TotalPointsAllowed,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		seasonStats = append(seasonStats, stats)
	}

	return seasonStats, nil
}

type ComparisonReport struct {
	HistoricalPlayer      string `json:"historical_player"`
	YearStart             int    `json:"year_start"`
	YearEnd               int    `json:"year_end"`
	HistoricalCompletions int    `json:"historical_completions"`
	CurrentPlayer         string `json:"current_player"`
	CurrentCompletions    int    `json:"current_completions"`
}

func GetComparisonWithHistoricalRecords() ([]ComparisonReport, error) {
	query := `
        SELECT h.player_name AS historical_player, h.year_start, h.year_end, 
               COALESCE(h.completions, 0) AS historical_completions,
               COALESCE(c.name, 'N/A') AS current_player, COALESCE(SUM(g.completions), 0) AS current_completions
        FROM historicalrecords h
        LEFT JOIN players c ON c.name = h.player_name
        LEFT JOIN playergamestats g ON g.player_id = c.player_id
        GROUP BY h.player_name, h.year_start, h.year_end, h.completions, c.name
        ORDER BY current_completions DESC;
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var reports []ComparisonReport
	for rows.Next() {
		var report ComparisonReport
		err := rows.Scan(&report.HistoricalPlayer, &report.YearStart, &report.YearEnd, &report.HistoricalCompletions,
			&report.CurrentPlayer, &report.CurrentCompletions)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

type PlayerStatsReport struct {
	PlayerName           string  `json:"player_name"`
	Position             string  `json:"position"`
	PassingYards         int     `json:"passing_yards"`
	PassingTDs           int     `json:"passing_tds"`
	Interceptions        int     `json:"interceptions"`
	RushingYards         int     `json:"rushing_yards"`
	RushingTDs           int     `json:"rushing_tds"`
	ReceivingYards       int     `json:"receiving_yards"`
	ReceivingTDs         int     `json:"receiving_tds"`
	Completions          int     `json:"completions"`
	PassAttempts         int     `json:"pass_attempts"`
	RushAttempts         int     `json:"rush_attempts"`
	Receptions           int     `json:"receptions"`
	CompletionPercentage float64 `json:"completion_percentage"`
	YardsPerAttempt      float64 `json:"yards_per_attempt"`
	QBRating             float64 `json:"qb_rating"`
	YardsPerCarry        float64 `json:"yards_per_carry"`
	ScrimmageYards       int     `json:"scrimmage_yards"`
	YardsPerScrimmage    float64 `json:"yards_per_scrimmage"`
	YardsPerReception    float64 `json:"yards_per_reception"`
}

func GetPlayerStatsByPosition(position string) ([]PlayerStatsReport, error) {
	query := `
        SELECT name, position,
            SUM(passing_yards) AS passing_yards,
            SUM(passing_tds) AS passing_tds,
            SUM(interceptions) AS interceptions,
            SUM(rushing_yards) AS rushing_yards,
            SUM(rushing_tds) AS rushing_tds,
            SUM(receiving_yards) AS receiving_yards,
            SUM(receiving_tds) AS receiving_tds,
            SUM(completions) AS completions,
            SUM(pass_attempts) AS pass_attempts,
            SUM(rush_attempts) AS rush_attempts,
            SUM(receptions) AS receptions,
            -- Estatísticas avançadas com COALESCE para evitar NULLs
            COALESCE((SUM(completions) / NULLIF(SUM(pass_attempts), 0)) * 100, 0) AS completion_percentage,
            COALESCE((SUM(passing_yards) / NULLIF(SUM(pass_attempts), 0)), 0) AS yards_per_attempt,
            COALESCE(((8.4 * SUM(passing_yards)) + (330 * SUM(passing_tds)) + (100 * SUM(completions)) - (200 * SUM(interceptions))) / NULLIF(SUM(pass_attempts), 0), 0) AS qb_rating,
            COALESCE((SUM(rushing_yards) / NULLIF(SUM(rush_attempts), 0)), 0) AS yards_per_carry,
            COALESCE((SUM(rushing_yards) + SUM(receiving_yards)), 0) AS scrimmage_yards,
            COALESCE((SUM(rushing_yards) + SUM(receiving_yards)) / NULLIF((SUM(rush_attempts) + SUM(receptions)), 0), 0) AS yards_per_scrimmage,
            COALESCE((SUM(receiving_yards) / NULLIF(SUM(receptions), 0)), 0) AS yards_per_reception
        FROM players
        JOIN playergamestats ON players.player_id = playergamestats.player_id
        WHERE position = ?
        GROUP BY name, position
        ORDER BY name;
    `

	fmt.Println("Executando consulta SQL para obter estatísticas avançadas dos jogadores")
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
			&report.PassingYards, &report.PassingTDs, &report.Interceptions,
			&report.RushingYards, &report.RushingTDs,
			&report.ReceivingYards, &report.ReceivingTDs,
			&report.Completions, &report.PassAttempts,
			&report.RushAttempts, &report.Receptions,
			&report.CompletionPercentage, &report.YardsPerAttempt, &report.QBRating,
			&report.YardsPerCarry, &report.ScrimmageYards, &report.YardsPerScrimmage,
			&report.YardsPerReception,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	fmt.Println("Consulta executada com sucesso e resultados processados")
	return reports, nil
}

type PlayerAverageStats struct {
	AvgCompletions    float64 `json:"avg_completions"`
	AvgPassingYards   float64 `json:"avg_passing_yards"`
	AvgPassingTDs     float64 `json:"avg_passing_tds"`
	AvgRushingYards   float64 `json:"avg_rushing_yards"`
	AvgRushingTDs     float64 `json:"avg_rushing_tds"`
	AvgReceivingYards float64 `json:"avg_receiving_yards"`
	AvgReceivingTDs   float64 `json:"avg_receiving_tds"`
}

func GetPlayerAverageStats(playerID int) (PlayerAverageStats, error) {
	query := `
        SELECT 
            AVG(completions) AS avg_completions,
            AVG(passing_yards) AS avg_passing_yards,
            AVG(passing_tds) AS avg_passing_tds,
            AVG(rushing_yards) AS avg_rushing_yards,
            AVG(rushing_tds) AS avg_rushing_tds,
            AVG(receiving_yards) AS avg_receiving_yards,
            AVG(receiving_tds) AS avg_receiving_tds
        FROM playergamestats
        WHERE player_id = ?
    `

	row := database.DB.QueryRow(query, playerID)
	var stats PlayerAverageStats
	err := row.Scan(
		&stats.AvgCompletions,
		&stats.AvgPassingYards,
		&stats.AvgPassingTDs,
		&stats.AvgRushingYards,
		&stats.AvgRushingTDs,
		&stats.AvgReceivingYards,
		&stats.AvgReceivingTDs,
	)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return stats, err
	}

	return stats, nil
}

type PredictionReport struct {
	PlayerID                int `json:"player_id"`
	PredictedCompletions    int `json:"predicted_completions"`
	RecordCompletions       int `json:"record_completions"`
	PredictedPassingYards   int `json:"predicted_passing_yards"`
	RecordPassingYards      int `json:"record_passing_yards"`
	PredictedRushingYards   int `json:"predicted_rushing_yards"`
	RecordRushingYards      int `json:"record_rushing_yards"`
	PredictedReceivingYards int `json:"predicted_receiving_yards"`
	RecordReceivingYards    int `json:"record_receiving_yards"`
}

func PredictRecordBreak(playerID int, seasonsRemaining int) (PredictionReport, error) {
	// Obter as médias anuais
	avgStats, err := GetPlayerAverageStats(playerID)
	if err != nil {
		return PredictionReport{}, err
	}

	// Obter os recordes máximos
	careerRecords, err := GetCareerRecords()
	if err != nil {
		return PredictionReport{}, err
	}

	prediction := PredictionReport{
		PlayerID:                playerID,
		PredictedCompletions:    int(avgStats.AvgCompletions * float64(seasonsRemaining)),
		RecordCompletions:       careerRecords.MaxCompletions,
		PredictedPassingYards:   int(avgStats.AvgPassingYards * float64(seasonsRemaining)),
		RecordPassingYards:      careerRecords.MaxPassingYards,
		PredictedRushingYards:   int(avgStats.AvgRushingYards * float64(seasonsRemaining)),
		RecordRushingYards:      careerRecords.MaxRushingYards,
		PredictedReceivingYards: int(avgStats.AvgReceivingYards * float64(seasonsRemaining)),
		RecordReceivingYards:    careerRecords.MaxReceivingYards,
	}

	return prediction, nil
}

type CareerRecords struct {
	MaxCompletions    int `json:"max_completions"`
	MaxPassingYards   int `json:"max_passing_yards"`
	MaxPassingTDs     int `json:"max_passing_tds"`
	MaxRushingYards   int `json:"max_rushing_yards"`
	MaxRushingTDs     int `json:"max_rushing_tds"`
	MaxReceivingYards int `json:"max_receiving_yards"`
	MaxReceivingTDs   int `json:"max_receiving_tds"`
}

func GetCareerRecords() (CareerRecords, error) {
	query := `
        SELECT 
            MAX(completions) AS max_completions,
            MAX(passing_yards) AS max_passing_yards,
            MAX(touchdowns) AS max_passing_tds,
            MAX(rush_yards) AS max_rushing_yards,
            MAX(rush_tds) AS max_rushing_tds,
            MAX(receiving_yards) AS max_receiving_yards,
            MAX(receiving_tds) AS max_receiving_tds
        FROM historicalrecords;
    `

	row := database.DB.QueryRow(query)
	var records CareerRecords
	err := row.Scan(
		&records.MaxCompletions,
		&records.MaxPassingYards,
		&records.MaxPassingTDs,
		&records.MaxRushingYards,
		&records.MaxRushingTDs,
		&records.MaxReceivingYards,
		&records.MaxReceivingTDs,
	)
	if err != nil {
		return records, err
	}

	return records, nil
}

type PlayerCareerStats struct {
	PlayerName           string `json:"player_name"`
	CareerCompletions    int    `json:"career_completions"`
	CareerPassingYards   int    `json:"career_passing_yards"`
	CareerPassingTDs     int    `json:"career_passing_tds"`
	CareerRushingYards   int    `json:"career_rushing_yards"`
	CareerRushingTDs     int    `json:"career_rushing_tds"`
	CareerReceivingYards int    `json:"career_receiving_yards"`
	CareerReceivingTDs   int    `json:"career_receiving_tds"`
}

func GetCurrentPlayerCareerStats() ([]PlayerCareerStats, error) {
	query := `
        SELECT 
            p.name AS player_name,
            SUM(g.completions) AS career_completions,
            SUM(g.passing_yards) AS career_passing_yards,
            SUM(g.passing_tds) AS career_passing_tds,
            SUM(g.rushing_yards) AS career_rushing_yards,
            SUM(g.rushing_tds) AS career_rushing_tds,
            SUM(g.receiving_yards) AS career_receiving_yards,
            SUM(g.receiving_tds) AS career_receiving_tds
        FROM players p
        LEFT JOIN playergamestats g ON g.player_id = p.player_id
        GROUP BY p.name
        ORDER BY career_passing_yards DESC;
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var playerStats []PlayerCareerStats
	for rows.Next() {
		var stats PlayerCareerStats
		err := rows.Scan(
			&stats.PlayerName,
			&stats.CareerCompletions,
			&stats.CareerPassingYards,
			&stats.CareerPassingTDs,
			&stats.CareerRushingYards,
			&stats.CareerRushingTDs,
			&stats.CareerReceivingYards,
			&stats.CareerReceivingTDs,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		playerStats = append(playerStats, stats)
	}

	return playerStats, nil
}

type ComparisonWithRecord struct {
	PlayerName           string `json:"player_name"`
	CareerCompletions    int    `json:"career_completions"`
	RecordCompletions    int    `json:"record_completions"`
	CareerPassingYards   int    `json:"career_passing_yards"`
	RecordPassingYards   int    `json:"record_passing_yards"`
	CareerRushingYards   int    `json:"career_rushing_yards"`
	RecordRushingYards   int    `json:"record_rushing_yards"`
	CareerReceivingYards int    `json:"career_receiving_yards"`
	RecordReceivingYards int    `json:"record_receiving_yards"`
}

func ComparePlayerStatsWithRecords() ([]ComparisonWithRecord, error) {
	// Obtenha os recordes de carreira
	records, err := GetCareerRecords()
	if err != nil {
		return nil, err
	}

	// Obtenha as estatísticas de carreira dos jogadores atuais
	playerStats, err := GetCurrentPlayerCareerStats()
	if err != nil {
		return nil, err
	}

	var comparisons []ComparisonWithRecord
	for _, stats := range playerStats {
		comparison := ComparisonWithRecord{
			PlayerName:           stats.PlayerName,
			CareerCompletions:    stats.CareerCompletions,
			RecordCompletions:    records.MaxCompletions,
			CareerPassingYards:   stats.CareerPassingYards,
			RecordPassingYards:   records.MaxPassingYards,
			CareerRushingYards:   stats.CareerRushingYards,
			RecordRushingYards:   records.MaxRushingYards,
			CareerReceivingYards: stats.CareerReceivingYards,
			RecordReceivingYards: records.MaxReceivingYards,
			// Adicione outras comparações aqui
		}
		comparisons = append(comparisons, comparison)
	}

	return comparisons, nil
}

type PlayerYearlyStats struct {
	Year           int `json:"year"`
	Completions    int `json:"completions"`
	PassingYards   int `json:"passing_yards"`
	PassingTDs     int `json:"passing_tds"`
	RushingYards   int `json:"rushing_yards"`
	RushingTDs     int `json:"rushing_tds"`
	ReceivingYards int `json:"receiving_yards"`
	ReceivingTDs   int `json:"receiving_tds"`
}

func GetPlayerCareerProgression(playerID int) ([]PlayerYearlyStats, error) {
	query := `
        SELECT 
            s.year,  -- Obtém o ano da tabela schedule
            SUM(g.completions) AS completions,
            SUM(g.passing_yards) AS passing_yards,
            SUM(g.passing_tds) AS passing_tds,
            SUM(g.rushing_yards) AS rushing_yards,
            SUM(g.rushing_tds) AS rushing_tds,
            SUM(g.receiving_yards) AS receiving_yards,
            SUM(g.receiving_tds) AS receiving_tds
        FROM playergamestats g
        JOIN schedule s ON g.schedule_id = s.id
        WHERE g.player_id = ?
        GROUP BY s.year
        ORDER BY s.year;
    `

	rows, err := database.DB.Query(query, playerID)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var yearlyStats []PlayerYearlyStats
	for rows.Next() {
		var stats PlayerYearlyStats
		err := rows.Scan(
			&stats.Year,
			&stats.Completions,
			&stats.PassingYards,
			&stats.PassingTDs,
			&stats.RushingYards,
			&stats.RushingTDs,
			&stats.ReceivingYards,
			&stats.ReceivingTDs,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		yearlyStats = append(yearlyStats, stats)
	}

	return yearlyStats, nil
}

type TopPlayerStats struct {
	PlayerName string `json:"player_name"`
	StatValue  int    `json:"stat_value"`
}

func GetTopPlayersBySeason(year int, category string) ([]TopPlayerStats, error) {
	query := fmt.Sprintf(`
        SELECT 
            p.name AS player_name,
            SUM(g.%s) AS stat_value
        FROM players p
        LEFT JOIN playergamestats g ON g.player_id = p.player_id
        JOIN schedule s ON g.schedule_id = s.id  -- Obtém o ano de schedule
        WHERE s.year = ?  -- Filtra pelo ano na tabela schedule
        GROUP BY p.name
        ORDER BY stat_value DESC
        LIMIT 10;
    `, category)

	rows, err := database.DB.Query(query, year)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var topPlayers []TopPlayerStats
	for rows.Next() {
		var stats TopPlayerStats
		err := rows.Scan(&stats.PlayerName, &stats.StatValue)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		topPlayers = append(topPlayers, stats)
	}

	return topPlayers, nil
}
