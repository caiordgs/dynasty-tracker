package services

import (
	"database/sql"
	"fmt"
	"log"
)

// Estrutura para a atribuição de time e cargo
type TeamAssignment struct {
	TeamID int    `json:"team_id"`
	Role   string `json:"role"`
	Year   int    `json:"year"`
}

// Função para atribuir time e cargo ao técnico
func AssignTeam(db *sql.DB, teamAssignment TeamAssignment) error {
	// Definindo a query para inserir a atribuição
	query := `INSERT INTO team_assignments (team_id, role, year) VALUES (?, ?, ?)`

	// Executando a query
	_, err := db.Exec(query, teamAssignment.TeamID, teamAssignment.Role, teamAssignment.Year)
	if err != nil {
		log.Println("Erro ao inserir a atribuição:", err)
		return fmt.Errorf("erro ao salvar atribuição: %v", err)
	}

	return nil
}
