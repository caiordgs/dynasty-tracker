package services

import (
	"database/sql"
	"dynastyTracker/database"
	"dynastyTracker/models"
)

// AddHistoricalRecord adiciona um novo recorde histórico ao banco de dados
func AddHistoricalRecord(record models.HistoricalRecord) error {
	_, err := database.DB.Exec(`INSERT INTO historicalrecords (school, player_name, year_start, year_end, completions, 
        attempts, completion_percentage, passing_yards, yards_per_attempt, touchdowns, interceptions, passer_rating,
        rush_attempts, rush_yards, yards_per_carry, rush_tds, receptions, receiving_yards, yards_per_catch, 
        receiving_tds, plays_from_scrimmage, yards_from_scrimmage, avg_yards_per_play, scrimmage_tds)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.School, record.PlayerName, record.YearStart, record.YearEnd, record.Completions, record.Attempts,
		record.CompletionPercentage, record.PassingYards, record.YardsPerAttempt, record.Touchdowns, record.Interceptions,
		record.PasserRating, record.RushAttempts, record.RushYards, record.YardsPerCarry, record.RushTDs,
		record.Receptions, record.ReceivingYards, record.YardsPerCatch, record.ReceivingTDs, record.PlaysFromScrimmage,
		record.YardsFromScrimmage, record.AvgYardsPerPlay, record.ScrimmageTDs)
	return err
}

// GetHistoricalRecord obtém um recorde histórico específico pelo ID
func GetHistoricalRecord(id int) (models.HistoricalRecord, error) {
	var record models.HistoricalRecord
	row := database.DB.QueryRow("SELECT * FROM historicalrecords WHERE record_id = ?", id)
	err := row.Scan(
		&record.RecordID, &record.School, &record.PlayerName, &record.YearStart, &record.YearEnd,
		&record.Completions, &record.Attempts, &record.CompletionPercentage, &record.PassingYards,
		&record.YardsPerAttempt, &record.Touchdowns, &record.Interceptions, &record.PasserRating,
		&record.RushAttempts, &record.RushYards, &record.YardsPerCarry, &record.RushTDs,
		&record.Receptions, &record.ReceivingYards, &record.YardsPerCatch, &record.ReceivingTDs,
		&record.PlaysFromScrimmage, &record.YardsFromScrimmage, &record.AvgYardsPerPlay, &record.ScrimmageTDs,
	)
	if err == sql.ErrNoRows {
		return record, err
	}
	return record, nil
}

// DeleteHistoricalRecord exclui um recorde histórico pelo ID
func DeleteHistoricalRecord(id int) error {
	_, err := database.DB.Exec("DELETE FROM historicalrecords WHERE record_id = ?", id)
	return err
}

// GetHistoricalRecords retorna todos os recordes históricos
func GetHistoricalRecords() ([]models.HistoricalRecord, error) {
	rows, err := database.DB.Query("SELECT * FROM historicalrecords")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.HistoricalRecord
	for rows.Next() {
		var record models.HistoricalRecord
		err := rows.Scan(
			&record.RecordID, &record.School, &record.PlayerName, &record.YearStart, &record.YearEnd,
			&record.Completions, &record.Attempts, &record.CompletionPercentage, &record.PassingYards,
			&record.YardsPerAttempt, &record.Touchdowns, &record.Interceptions, &record.PasserRating,
			&record.RushAttempts, &record.RushYards, &record.YardsPerCarry, &record.RushTDs,
			&record.Receptions, &record.ReceivingYards, &record.YardsPerCatch, &record.ReceivingTDs,
			&record.PlaysFromScrimmage, &record.YardsFromScrimmage, &record.AvgYardsPerPlay, &record.ScrimmageTDs,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func UpdateHistoricalRecord(record models.HistoricalRecord) error {
	_, err := database.DB.Exec(`UPDATE historicalrecords SET school=?, player_name=?, year_start=?, year_end=?, 
        completions=?, attempts=?, completion_percentage=?, passing_yards=?, yards_per_attempt=?, touchdowns=?, 
        interceptions=?, passer_rating=?, rush_attempts=?, rush_yards=?, yards_per_carry=?, rush_tds=?, 
        receptions=?, receiving_yards=?, yards_per_catch=?, receiving_tds=?, plays_from_scrimmage=?, 
        yards_from_scrimmage=?, avg_yards_per_play=?, scrimmage_tds=? WHERE record_id=?`,
		record.School, record.PlayerName, record.YearStart, record.YearEnd, record.Completions, record.Attempts,
		record.CompletionPercentage, record.PassingYards, record.YardsPerAttempt, record.Touchdowns, record.Interceptions,
		record.PasserRating, record.RushAttempts, record.RushYards, record.YardsPerCarry, record.RushTDs,
		record.Receptions, record.ReceivingYards, record.YardsPerCatch, record.ReceivingTDs, record.PlaysFromScrimmage,
		record.YardsFromScrimmage, record.AvgYardsPerPlay, record.ScrimmageTDs, record.RecordID)
	return err
}
