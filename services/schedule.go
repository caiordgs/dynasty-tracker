package services

import (
	"database/sql"
	"dynastyTracker/database"
	"dynastyTracker/models"
	"fmt"
)

// GetSchedules retorna todos os jogos do calendário
func GetSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule

	query := "SELECT id, team_name, year, week, opponent, team_ranking, opponent_ranking, team_points, opponent_points, result, site FROM schedule"
	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Println("Erro na consulta SQL para obter o calendário:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.TeamName, &schedule.Year, &schedule.Week, &schedule.Opponent, &schedule.TeamRanking, &schedule.OpponentRanking, &schedule.TeamPoints, &schedule.OpponentPoints, &schedule.Result, &schedule.Site); err != nil {
			fmt.Println("Erro ao escanear resultados:", err)
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// AddSchedule adiciona um novo jogo ao calendário
func AddSchedule(schedule models.Schedule) error {
	_, err := database.DB.Exec(`INSERT INTO schedule (team_id, team_name, year, week, opponent, team_ranking,
        opponent_ranking, team_points, opponent_points, result, site)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		schedule.TeamID, schedule.TeamName, schedule.Year, schedule.Week, schedule.Opponent, schedule.TeamRanking,
		schedule.OpponentRanking, schedule.TeamPoints, schedule.OpponentPoints, schedule.Result, schedule.Site)
	return err
}

// GetSchedule obtém um jogo específico do calendário pelo ID
func GetSchedule(id int) (models.Schedule, error) {
	var schedule models.Schedule
	row := database.DB.QueryRow("SELECT * FROM schedule WHERE id = ?", id)
	err := row.Scan(
		&schedule.ID, &schedule.TeamID, &schedule.TeamName, &schedule.Year, &schedule.Week,
		&schedule.Opponent, &schedule.TeamRanking, &schedule.OpponentRanking, &schedule.TeamPoints,
		&schedule.OpponentPoints, &schedule.Result, &schedule.Site,
	)
	if err == sql.ErrNoRows {
		return schedule, err
	}
	return schedule, nil
}

// DeleteSchedule exclui um jogo específico do calendário pelo ID
func DeleteSchedule(id int) error {
	_, err := database.DB.Exec("DELETE FROM schedule WHERE id = ?", id)
	return err
}

func UpdateSchedule(schedule models.Schedule) error {
	_, err := database.DB.Exec(`UPDATE schedule SET team_id=?, team_name=?, year=?, week=?, opponent=?, team_ranking=?, 
        opponent_ranking=?, team_points=?, opponent_points=?, result=?, site=? WHERE id=?`,
		schedule.TeamID, schedule.TeamName, schedule.Year, schedule.Week, schedule.Opponent, schedule.TeamRanking,
		schedule.OpponentRanking, schedule.TeamPoints, schedule.OpponentPoints, schedule.Result, schedule.Site, schedule.ID)
	return err
}

func GetSchedulesWithFilters(year int, week int) ([]models.Schedule, error) {
	query := "SELECT * FROM schedule WHERE 1=1"
	var args []interface{}

	if year > 0 {
		query += " AND year = ?"
		args = append(args, year)
	}
	if week > 0 {
		query += " AND week = ?"
		args = append(args, week)
	}

	fmt.Printf("Query: %s, Args: %v\n", query, args) // Log temporário para depuração da query

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err) // Log detalhado do erro
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ID, &schedule.TeamName, &schedule.TeamID, &schedule.Week, &schedule.Opponent,
			&schedule.TeamRanking, &schedule.OpponentRanking, &schedule.TeamPoints, &schedule.OpponentPoints,
			&schedule.Result, &schedule.Site, &schedule.Year,
		)
		if err != nil {
			fmt.Printf("Erro ao escanear resultados: %v\n", err)
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
