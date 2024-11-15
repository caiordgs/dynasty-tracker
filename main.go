package main

import (
	"dynastyTracker/database"
	"dynastyTracker/models"
	"dynastyTracker/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Função principal
func main() {
	// Inicializar a conexão com o banco de dados
	err := database.InitDB()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Definir rotas para cada recurso
	// Jogadores
	http.HandleFunc("/api/players", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w, r) // Sem o ponteiro, passando diretamente o http.ResponseWriter
		playersHandler(w, r)
	})
	http.HandleFunc("/api/players/", playerHandler) // Busca jogador por ID

	// Calendário
	http.HandleFunc("/api/schedule", scheduleHandler)      // Lista e adiciona jogos no calendário
	http.HandleFunc("/api/schedule/", scheduleItemHandler) // Busca jogo do calendário por ID

	// Recordes Históricos
	http.HandleFunc("/api/records", recordsHandler) // Lista e adiciona recordes históricos
	http.HandleFunc("/api/records/", recordHandler) // Busca recorde histórico por ID

	// Consultas Específicas
	http.HandleFunc("/api/schedule/search", scheduleSearchHandler)
	http.HandleFunc("/api/players/search", playerSearchHandler)
	http.HandleFunc("/api/records/search", recordSearchHandler)

	// Relatórios
	http.HandleFunc("/api/reports/team-performance", teamPerformanceHandler)
	http.HandleFunc("/api/reports/player-stats", playerStatsHandler)
	http.HandleFunc("/api/reports/season-summary", seasonSummaryHandler)
	http.HandleFunc("/api/reports/comparison-records", comparisonWithHistoricalRecordsHandler)
	http.HandleFunc("/api/reports/player-records-comparison", playerRecordsComparisonHandler)
	http.HandleFunc("/api/reports/player-career-progression", playerCareerProgressionHandler)
	http.HandleFunc("/api/reports/top-players", topPlayersBySeasonHandler)
	http.HandleFunc("/api/reports/team-season-comparison", teamSeasonComparisonHandler)
	http.HandleFunc("/api/reports/record-break-prediction", recordBreakPredictionHandler)

	// Recrutas
	http.HandleFunc("/api/recruits/add", addRecruitHandler)
	http.HandleFunc("/api/players/add", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w, r) // Sem o ponteiro, passando diretamente o http.ResponseWriter
		addPlayerHandler(w, r)
	})
	http.HandleFunc("/api/teams", teamsHandler)             // Para acessar os times
	http.HandleFunc("/api/teams/assign", assignTeamHandler) // Para atribuir um time a um técnico

	// Iniciar o servidor
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Caso a requisição seja do tipo OPTIONS, retorna imediatamente
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// Handler para listar ou adicionar jogadores
// playersHandler no backend
func playersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		players, err := services.GetPlayers()
		if err != nil {
			http.Error(w, "Erro ao obter jogadores", http.StatusInternalServerError)
			return
		}

		// Log para ver os dados antes de enviar
		fmt.Println("Jogadores encontrados:", players)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(players)
	case http.MethodPost:
		var player models.Player
		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, "Erro ao decodificar jogador", http.StatusBadRequest)
			return
		}
		err = services.AddPlayer(player)
		if err != nil {
			http.Error(w, "Erro ao adicionar jogador", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Jogador adicionado com sucesso"})
	}
}

func playerHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		player, err := services.GetPlayer(id)
		if err != nil {
			http.Error(w, "Jogador não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)

	case http.MethodPut: // Adicionando a funcionalidade PUT
		var player models.Player
		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, "Erro ao decodificar jogador", http.StatusBadRequest)
			return
		}
		player.PlayerID = id // Certifique-se de usar o ID correto
		err = services.UpdatePlayer(player)
		if err != nil {
			http.Error(w, "Erro ao atualizar jogador", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Jogador atualizado com sucesso"})
	}
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		schedules, err := services.GetSchedules()
		if err != nil {
			fmt.Println("Erro ao obter calendário:", err) // Log detalhado
			http.Error(w, "Erro ao obter calendário", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schedules)
	}
}

func scheduleItemHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		schedule, err := services.GetSchedule(id)
		if err != nil {
			http.Error(w, "Jogo não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schedule)

	case http.MethodPut: // Adicionando a funcionalidade PUT
		var schedule models.Schedule
		err := json.NewDecoder(r.Body).Decode(&schedule)
		if err != nil {
			http.Error(w, "Erro ao decodificar jogo", http.StatusBadRequest)
			return
		}
		schedule.ID = id // Certifique-se de usar o ID correto
		err = services.UpdateSchedule(schedule)
		if err != nil {
			http.Error(w, "Erro ao atualizar jogo", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Jogo atualizado com sucesso"})
	}
}

func recordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		records, err := services.GetHistoricalRecords()
		if err != nil {
			http.Error(w, "Erro ao obter recordes históricos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)

	case http.MethodPost:
		var record models.HistoricalRecord
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			fmt.Println("Erro ao decodificar JSON do recorde:", err) // Log do erro
			http.Error(w, "Erro ao decodificar recorde", http.StatusBadRequest)
			return
		}
		err = services.AddHistoricalRecord(record)
		if err != nil {
			fmt.Println("Erro ao adicionar recorde no banco de dados:", err) // Log do erro
			http.Error(w, "Erro ao adicionar recorde", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Recorde adicionado com sucesso"})
	}
}

func recordHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		record, err := services.GetHistoricalRecord(id)
		if err != nil {
			http.Error(w, "Recorde não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(record)

	case http.MethodPut: // Adicionando a funcionalidade PUT
		var record models.HistoricalRecord
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			http.Error(w, "Erro ao decodificar recorde", http.StatusBadRequest)
			return
		}
		record.RecordID = id // Certifique-se de usar o ID correto
		err = services.UpdateHistoricalRecord(record)
		if err != nil {
			http.Error(w, "Erro ao atualizar recorde", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Recorde atualizado com sucesso"})
	}
}

func extractID(path string) int {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0 // Retorna 0 caso não consiga converter o ID
	}
	return id
}

func scheduleSearchHandler(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	weekParam := r.URL.Query().Get("week")

	var year, week int
	var err error

	if yearParam != "" {
		year, err = strconv.Atoi(yearParam)
		if err != nil {
			http.Error(w, "Parâmetro 'year' inválido", http.StatusBadRequest)
			return
		}
	}

	if weekParam != "" {
		week, err = strconv.Atoi(weekParam)
		if err != nil {
			http.Error(w, "Parâmetro 'week' inválido", http.StatusBadRequest)
			return
		}
	}

	schedules, err := services.GetSchedulesWithFilters(year, week)
	if err != nil {
		http.Error(w, "Erro ao buscar jogos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)

}

func playerSearchHandler(w http.ResponseWriter, r *http.Request) {
	position := r.URL.Query().Get("position")
	teamIDParam := r.URL.Query().Get("team_id")

	var teamID int
	if teamIDParam != "" {
		teamID, _ = strconv.Atoi(teamIDParam)
	}

	players, err := services.GetPlayersWithFilters(position, teamID)
	if err != nil {
		http.Error(w, "Erro ao buscar jogadores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func recordSearchHandler(w http.ResponseWriter, r *http.Request) {
	school := r.URL.Query().Get("school")
	playerName := r.URL.Query().Get("player_name")

	records, err := services.GetHistoricalRecordsWithFilters(school, playerName)
	if err != nil {
		http.Error(w, "Erro ao buscar recordes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func teamPerformanceHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := services.GetTeamPerformanceBySeason()
	if err != nil {
		http.Error(w, "Erro ao obter desempenho do time", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}

func playerStatsHandler(w http.ResponseWriter, r *http.Request) {
	position := r.URL.Query().Get("position")
	if position == "" {
		http.Error(w, "Parâmetro 'position' é obrigatório", http.StatusBadRequest)
		return
	}

	reports, err := services.GetPlayerStatsByPosition(position)
	if err != nil {
		http.Error(w, "Erro ao gerar relatório de estatísticas dos jogadores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func seasonSummaryHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "Ano é necessário", http.StatusBadRequest)
		return
	}

	// Converte o ano de string para int
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Ano inválido", http.StatusBadRequest)
		return
	}

	reports, err := services.GetSeasonSummary(year)
	if err != nil {
		http.Error(w, "Erro ao obter resumo da temporada", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}

func comparisonWithHistoricalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := services.GetComparisonWithHistoricalRecords()
	if err != nil {
		http.Error(w, "Erro ao obter comparação com recordes históricos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}

func playerRecordsComparisonHandler(w http.ResponseWriter, r *http.Request) {
	comparisons, err := services.ComparePlayerStatsWithRecords()
	if err != nil {
		http.Error(w, "Erro ao comparar estatísticas dos jogadores com recordes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comparisons)
}

func playerCareerProgressionHandler(w http.ResponseWriter, r *http.Request) {
	playerIDStr := r.URL.Query().Get("player_id")
	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		http.Error(w, "ID do jogador inválido", http.StatusBadRequest)
		return
	}

	yearlyStats, err := services.GetPlayerCareerProgression(playerID)
	if err != nil {
		http.Error(w, "Erro ao obter evolução de carreira", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(yearlyStats)
}

func topPlayersBySeasonHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	category := r.URL.Query().Get("category")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Ano inválido", http.StatusBadRequest)
		return
	}

	topPlayers, err := services.GetTopPlayersBySeason(year, category)
	if err != nil {
		http.Error(w, "Erro ao obter ranking de jogadores", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(topPlayers)
}

func teamSeasonComparisonHandler(w http.ResponseWriter, r *http.Request) {
	teamIDStr := r.URL.Query().Get("team_id")
	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		http.Error(w, "ID do time inválido", http.StatusBadRequest)
		return
	}

	seasonStats, err := services.GetTeamSeasonComparison(teamID)
	if err != nil {
		http.Error(w, "Erro ao obter comparação de temporadas", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(seasonStats)
}

func recordBreakPredictionHandler(w http.ResponseWriter, r *http.Request) {
	playerIDStr := r.URL.Query().Get("player_id")
	seasonsRemainingStr := r.URL.Query().Get("seasons_remaining")

	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		http.Error(w, "ID do jogador inválido", http.StatusBadRequest)
		return
	}

	seasonsRemaining, err := strconv.Atoi(seasonsRemainingStr)
	if err != nil {
		http.Error(w, "Número de temporadas restantes inválido", http.StatusBadRequest)
		return
	}

	prediction, err := services.PredictRecordBreak(playerID, seasonsRemaining)
	if err != nil {
		http.Error(w, "Erro ao gerar predição de quebra de recorde", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prediction)
}

func addPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	// Decodificar o corpo da requisição JSON para a estrutura Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "Erro ao decodificar dados do jogador", http.StatusBadRequest)
		return
	}

	// Buscar o team_id com base no nome do time
	teamID, err := getTeamIDByName(player.TeamName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Atualizar o player com o team_id
	player.TeamID = teamID

	// Inserir o jogador no banco de dados
	err = services.AddPlayer(player)
	if err != nil {
		http.Error(w, "Erro ao adicionar jogador", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Jogador adicionado com sucesso"})
}

func addPlayerGameStatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats models.PlayerGameStats

	err := json.NewDecoder(r.Body).Decode(&stats)
	if err != nil {
		http.Error(w, "Erro ao decodificar estatísticas do jogo", http.StatusBadRequest)
		return
	}

	err = services.AddPlayerGameStats(stats)
	if err != nil {
		http.Error(w, "Erro ao adicionar estatísticas do jogo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Estatísticas do jogo adicionadas com sucesso"})
}

func addRecruitHandler(w http.ResponseWriter, r *http.Request) {
	var recruit models.Recruit

	// Decodificar JSON do corpo da requisição
	err := json.NewDecoder(r.Body).Decode(&recruit)
	if err != nil {
		http.Error(w, "Erro ao decodificar dados do recruta", http.StatusBadRequest)
		return
	}

	// Inserir recruta na tabela recruits
	err = services.AddRecruit(recruit)
	if err != nil {
		http.Error(w, "Erro ao adicionar recruta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recruta adicionado com sucesso"})
}

func addRecruitedPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var recruit models.Recruit

	// Decodificar JSON do corpo da requisição
	err := json.NewDecoder(r.Body).Decode(&recruit)
	if err != nil {
		http.Error(w, "Erro ao decodificar dados do recruta", http.StatusBadRequest)
		return
	}

	// Chamar a função de serviço para adicionar o recruta
	err = services.AddRecruit(recruit)
	if err != nil {
		http.Error(w, "Erro ao adicionar recruta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recruta adicionado com sucesso"})
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := services.GetTeams()
	if err != nil {
		http.Error(w, "Erro ao obter os times", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func assignTeamHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Obter dados de atribuição do técnico ao time
		var assignment models.TeamAssignment
		err := json.NewDecoder(r.Body).Decode(&assignment)
		if err != nil {
			http.Error(w, "Erro ao decodificar dados da atribuição", http.StatusBadRequest)
			return
		}

		// Inserir dados na tabela team_assignments
		err = services.AssignTeamToCoach(assignment)
		if err != nil {
			http.Error(w, "Erro ao atribuir time", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Atribuição realizada com sucesso!"})
	}
}
