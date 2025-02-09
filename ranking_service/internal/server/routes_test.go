package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestFetchRanking(t *testing.T) {
	mockRedis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer mockRedis.Close()

	mockRedis.ZAdd(context.Background(), "ranking", &redis.Z{
		Score:  1,
		Member: "video1",
	})
	mockRedis.ZAdd(context.Background(), "ranking", &redis.Z{
		Score:  2,
		Member: "video2",
	})

	// Tạo một request giả lập
	req, err := http.NewRequest("GET", "/fetch-ranking", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		var leaderboard, err = mockRedis.ZRevRangeWithScores(ctx, "ranking", 0, -1).Result()
		if err != nil {
			log.Printf("Error getting sorted ranking ascending: %v\n", err)
			http.Error(w, "Error fetching ranking", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(leaderboard)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(data)); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `[{"Score":200,"Member":"video1"},{"Score":100,"Member":"video2"}]`
	assert.JSONEq(t, expected, rr.Body.String())
}
