package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"user_interaction_service/internal"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func IncrementScore(video string, increment int64) {
	err := internal.RedisServer.ZIncrBy(ctx, "ranking", float64(increment), video).Err()
	if err != nil {
		log.Printf("Error incrementing score for %s: %v\n", video, err)
	}
}

func AddVideoScores(userID, videoID string, types map[string]int) {
	for t, score := range types {
		currentScore, err := internal.RedisServer.HGet(ctx, fmt.Sprintf("user:%s:videos", userID), videoID+"|"+t).Result()

		if err == redis.Nil {
			currentScore = "0"
		} else if err != nil {
			log.Fatalf("Error getting current score: %v", err)
		}

		currentScoreInt, err := strconv.Atoi(currentScore)
		if err != nil {
			log.Fatalf("Error converting current score to int: %v", err)
		}

		newScore := currentScoreInt + score

		err = internal.RedisServer.HSet(ctx, fmt.Sprintf("user:%s:videos", userID), videoID+"|"+t, newScore).Err()
		if err != nil {
			log.Fatalf("Error updating video score: %v", err)
		}
	}
}

func CalculateTotalScore(userID, videoID string) int {
	typesAndScores, err := internal.RedisServer.HGetAll(ctx, fmt.Sprintf("user:%s:videos", userID)).Result()
	if err != nil {
		log.Fatalf("Error getting video scores from Redis: %v", err)
	}

	totalScore := 0
	for key, value := range typesAndScores {
		if len(key) > len(videoID) && key[:len(videoID)] == videoID {
			typeKey := key[len(videoID)+1:]
			score, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Error converting score for type %s: %v", typeKey, err)
			}

			totalScore += score
		}
	}

	return totalScore
}

func GetTopVideos(userID string) []string {
	// Lấy danh sách video có điểm cao nhất từ Sorted Set
	videos, err := internal.RedisServer.ZRevRange(ctx, fmt.Sprintf("user:%s:video_ranking", userID), 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting top videos: %v", err)
	}
	return videos
}

func AddToRanking(userID, videoID string, score int) {
	err := internal.RedisServer.ZAdd(ctx, fmt.Sprintf("user:%s:video_ranking", userID), &redis.Z{
		Score:  float64(score),
		Member: videoID,
	}).Err()
	if err != nil {
		log.Fatalf("Error adding video to ranking: %v", err)
	}
}

func AddVideoToGlobalRanking(videoID string, score int) error {
	rankingKey := "global:video:ranking"

	_, err := internal.RedisServer.ZAdd(ctx, rankingKey, &redis.Z{
		Score:  float64(score),
		Member: videoID,
	}).Result()
	if err != nil {
		return fmt.Errorf("Error adding video to global ranking: %v", err)
	}
	return nil
}
