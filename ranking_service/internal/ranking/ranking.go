package ranking

import (
	"encoding/json"
	"log"
	"net/http"
	"ranking_service/internal/utils/redis"
)

func FetchRanking(w http.ResponseWriter, r *http.Request) {
	var leaderboard, err = redis.GetTopVideosFromGlobalRanking()
	if err != nil {
		log.Printf("Error getting sorted ranking ascending: %v\n", err)
		return
	}

	data, err := json.Marshal(leaderboard)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	if _, err := w.Write([]byte(data)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func FetchUserRanking(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "userID isn't empty", http.StatusBadRequest)
		return
	}

	var leaderboard = redis.GetTopVideos(userID)

	data, err := json.Marshal(leaderboard)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	if _, err := w.Write([]byte(data)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
