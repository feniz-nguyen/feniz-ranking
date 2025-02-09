package redis

import (
	"context"
	"fmt"
	"log"
	"ranking_service/internal"
)

type Video struct {
	ID    string
	Score int
}

var ctx = context.Background()

func GetTopVideos(userID string) []string {
	videos, err := internal.RedisServer.ZRevRange(ctx, fmt.Sprintf("user:%s:video_ranking", userID), 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting top videos: %v", err)
	}
	return videos
}

func GetTopVideosFromGlobalRanking() ([]Video, error) {
	rankingKey := "global:video:ranking"

	videoIDs, err := internal.RedisServer.ZRevRange(ctx, rankingKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("Error getting videos from global ranking: %v", err)
	}

	var videos []Video
	for _, videoID := range videoIDs {
		score, err := internal.RedisServer.ZScore(ctx, rankingKey, videoID).Result()
		if err != nil {
			return nil, fmt.Errorf("Error getting score for video %s: %v", videoID, err)
		}
		videos = append(videos, Video{ID: videoID, Score: int(score)})
	}

	return videos, nil
}
