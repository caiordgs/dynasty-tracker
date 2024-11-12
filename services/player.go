package services

import (
	"database/sql"
	"dynastyTracker/database"
	"dynastyTracker/models"
	"fmt"
)

// GetPlayers retorna a lista de todos os jogadores
func GetPlayers() ([]models.Player, error) {
	var players []models.Player

	query := "SELECT player_id, name, position, overall, class_year, team_id FROM players"
	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Println("Erro ao executar a consulta SQL:", err) // Log do erro SQL
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player models.Player
		if err := rows.Scan(&player.PlayerID, &player.Name, &player.Position, &player.Overall, &player.ClassYear, &player.TeamID); err != nil {
			fmt.Println("Erro ao escanear dados do jogador:", err) // Log de erro de escaneamento
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

type Player struct {
	Name      string `json:"name"`
	Position  string `json:"position"`
	Overall   int    `json:"overall"`
	ClassYear string `json:"class_year"`
	TeamID    int    `json:"team_id"`
}

// AddPlayer adiciona um novo jogador ao banco de dados
func AddPlayer(player models.Player) error {
	// Verificar o limite do elenco
	var playerCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM players WHERE team_id = ?", player.TeamID).Scan(&playerCount)
	if err != nil {
		fmt.Printf("Erro ao contar jogadores: %v\n", err)
		return err
	}

	if playerCount >= 85 {
		return fmt.Errorf("O elenco atingiu o limite máximo de 85 jogadores")
	}

	// Inserir o jogador se o limite não foi atingido
	query := `
        INSERT INTO players (name, position, overall, games_played, games_started, snaps_played, class_year, team_id, recruitment_source)
        VALUES (?, ?, ?, 0, 0, 0, ?, ?, ?);
    `
	_, err = database.DB.Exec(query, player.Name, player.Position, player.Overall, player.ClassYear, player.TeamID, player.RecruitmentSource)
	if err != nil {
		fmt.Printf("Erro ao adicionar jogador: %v\n", err)
		return err
	}
	return nil
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

func UpdatePlayer(player models.Player) error {
	_, err := database.DB.Exec(`UPDATE players SET name=?, position=?, overall=?, games_played=?, games_started=?, 
        snaps_played=?, class_year=?, team_id=? WHERE player_id=?`,
		player.Name, player.Position, player.Overall, player.GamesPlayed, player.GamesStarted,
		player.SnapsPlayed, player.ClassYear, player.TeamID, player.PlayerID)
	return err
}

func GetPlayersWithFilters(position string, teamID int) ([]models.Player, error) {
	query := "SELECT * FROM players WHERE 1=1"
	var args []interface{}

	if position != "" {
		query += " AND position = ?"
		args = append(args, position)
	}
	if teamID > 0 {
		query += " AND team_id = ?"
		args = append(args, teamID)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(&player.PlayerID, &player.Name, &player.Position, &player.Overall, &player.GamesPlayed,
			&player.GamesStarted, &player.SnapsPlayed, &player.ClassYear, &player.TeamID)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func PromoteRecruits(currentYear int) error {
	recruitmentYear := currentYear - 1

	rows, err := database.DB.Query(`
        SELECT recruit_id, player_name, position, overall, class, team_id, recruitment_source
        FROM recruits
        WHERE recruitment_year = ?
    `, recruitmentYear)
	if err != nil {
		fmt.Printf("Erro ao buscar recrutas: %v\n", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var recruit models.Recruit
		err := rows.Scan(&recruit.RecruitID, &recruit.PlayerName, &recruit.Position, &recruit.Overall, &recruit.Class, &recruit.TeamID, &recruit.RecruitmentSource)
		if err != nil {
			fmt.Printf("Erro ao escanear recruta: %v\n", err)
			return err
		}

		_, err = database.DB.Exec(`
            INSERT INTO players (name, position, overall, games_played, games_started, snaps_played, class_year, team_id, recruitment_source)
            VALUES (?, ?, ?, 0, 0, 0, ?, ?, ?);
        `, recruit.PlayerName, recruit.Position, recruit.Overall, recruit.Class, recruit.TeamID, recruit.RecruitmentSource)
		if err != nil {
			fmt.Printf("Erro ao promover recruta: %v\n", err)
			return err
		}
	}

	_, err = database.DB.Exec("DELETE FROM recruits WHERE recruitment_year = ?", recruitmentYear)
	if err != nil {
		fmt.Printf("Erro ao remover recrutas promovidos: %v\n", err)
		return err
	}

	return nil
}
