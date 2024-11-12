package services

import (
	"dynastyTracker/database"
	"dynastyTracker/models"
	"fmt"
)

// Função para obter todos os times
func GetTeams() ([]models.Team, error) {
	var teams []models.Team
	// Fazendo a consulta de todas as colunas
	rows, err := database.DB.Query("SELECT team_id, year, role, school, mascot, abbreviation, alt_name1, color, alt_color, logo_1, logo_2, twitter, location_venue_id, location_name, location_city, location_state FROM teams")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar times: %v", err)
	}
	defer rows.Close()

	// Iterando sobre os resultados da consulta e armazenando na lista de times
	for rows.Next() {
		var team models.Team
		// Preenchendo os campos de team com o resultado da consulta
		if err := rows.Scan(
			&team.TeamID,
			&team.Year,
			&team.Role,
			&team.School,
			&team.Mascot,
			&team.Abbreviation,
			&team.AltName1,
			&team.Color,
			&team.AltColor,
			&team.Logo1,
			&team.Logo2,
			&team.Twitter,
			&team.LocationVenueID,
			&team.LocationName,
			&team.LocationCity,
			&team.LocationState,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler dados do time: %v", err)
		}
		// Adicionando o time à lista
		teams = append(teams, team)
	}
	// Retornando a lista de times
	return teams, nil
}

// Função para atribuir um time a um técnico
func AssignTeamToCoach(assignment models.TeamAssignment) error {
	_, err := database.DB.Exec(`
		INSERT INTO team_assignments (team_id, coach_id, year, role)
		VALUES (?, ?, ?, ?)`,
		assignment.TeamID, assignment.CoachID, assignment.Year, assignment.Role)
	if err != nil {
		return fmt.Errorf("Erro ao atribuir time ao técnico: %v", err)
	}
	return nil
}
