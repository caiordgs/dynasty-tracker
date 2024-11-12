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

// Handler para operações com um jogador específico
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

	case http.MethodDelete:
		err := services.DeletePlayer(id)
		if err != nil {
			http.Error(w, "Erro ao excluir jogador", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
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

// Handler para operações com um jogo específico
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

	case http.MethodDelete:
		err := services.DeleteSchedule(id)
		if err != nil {
			http.Error(w, "Erro ao excluir jogo", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
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

// Handler para operações com um recorde específico
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

	case http.MethodDelete:
		err := services.DeleteHistoricalRecord(id)
		if err != nil {
			http.Error(w, "Erro ao excluir recorde", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
