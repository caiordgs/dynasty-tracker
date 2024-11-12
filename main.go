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
	http.HandleFunc("/api/players", playersHandler)
	http.HandleFunc("/api/players/", playerHandler) // Com ID específico
	http.HandleFunc("/api/schedule", scheduleHandler)
	http.HandleFunc("/api/schedule/", scheduleItemHandler) // Com ID específico
	http.HandleFunc("/api/records", recordsHandler)
	http.HandleFunc("/api/records/", recordHandler) // Com ID específico
	http.HandleFunc("/api/schedule/search", scheduleSearchHandler)
	http.HandleFunc("/api/players/search", playerSearchHandler)
	http.HandleFunc("/api/records/search", recordSearchHandler)

	// Iniciar o servidor
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler para listar ou adicionar jogadores
func playersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		players, err := services.GetPlayers()
		if err != nil {
			http.Error(w, "Erro ao obter jogadores", http.StatusInternalServerError)
			return
		}
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
			http.Error(w, "Erro ao obter calendário", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schedules)

	case http.MethodPost:
		var schedule models.Schedule
		err := json.NewDecoder(r.Body).Decode(&schedule)
		if err != nil {
			fmt.Println("Erro ao decodificar JSON do jogo:", err) // Log do erro
			http.Error(w, "Erro ao decodificar jogo", http.StatusBadRequest)
			return
		}
		err = services.AddSchedule(schedule)
		if err != nil {
			fmt.Println("Erro ao adicionar jogo no banco de dados:", err) // Log do erro
			http.Error(w, "Erro ao adicionar jogo", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Jogo adicionado com sucesso"})
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

// Função para extrair o ID do URL
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
