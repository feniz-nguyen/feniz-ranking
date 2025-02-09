package interaction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"user_interaction_service/internal/utils/kafka"
)

type Interaction struct {
	UserID string `json:"user_id"`
	Title  string `json:"video_name"`
	Type   string `json:"type"`
}

func InteractionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var inter Interaction
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inter)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	_, err = kafka.SendMessage(inter.Title, fmt.Sprintf("%v:%v", inter.UserID, inter.Type), "video-interactions")
	if err != nil {
		http.Error(w, "Kafka can't send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
