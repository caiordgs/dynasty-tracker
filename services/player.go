package services

import (
	"database/sql"
	"dynastyTracker/database"
	"dynastyTracker/models"
)

// GetPlayers retorna a lista de todos os jogadores
func GetPlayers() ([]models.Player, error) {
	rows, err := database.DB.Query("SELECT * FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(&player.PlayerID, &player.Name, &player.Position, &player.Overall,
			&player.GamesPlayed, &player.GamesStarted, &player.SnapsPlayed, &player.ClassYear, &player.TeamID)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

// AddPlayer adiciona um novo jogador ao banco de dados
func AddPlayer(player models.Player) error {
	_, err := database.DB.Exec(`INSERT INTO players (name, position, overall, games_played, games_started, snaps_played, class_year, team_id)
                                 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		player.Name, player.Position, player.Overall, player.GamesPlayed, player.GamesStarted, player.SnapsPlayed, player.ClassYear, player.TeamID)
	return err
}

// GetPlayer obtém um jogador específico pelo ID
func GetPlayer(id int) (models.Player, error) {
	var player models.Player
	row := database.DB.QueryRow("SELECT * FROM players WHERE player_id = ?", id)
	err := row.Scan(&player.PlayerID, &player.Name, &player.Position, &player.Overall,
		&player.GamesPlayed, &player.GamesStarted, &player.SnapsPlayed, &player.ClassYear, &player.TeamID)
	if err == sql.ErrNoRows {
		return player, err
	}
	return player, nil
}

// DeletePlayer exclui um jogador pelo ID
func DeletePlayer(id int) error {
	_, err := database.DB.Exec("DELETE FROM players WHERE player_id = ?", id)
	return err
}
