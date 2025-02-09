package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"
	"user_interaction_service/internal/utils/redis"

	"github.com/IBM/sarama"
)

func Consumer(topic string) {
	fmt.Println("Go to consumer")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create a new consumer
	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_URL")}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v\n", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting partition consumer: %v\n", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Println(msg)
			key := string(msg.Key)
			value := string(msg.Value)

			keys := strings.Split(key, ":")
			if len(keys) == 2 {
				redis.AddVideoScores(keys[0], value, map[string]int{
					keys[1]: 1,
				})

				score := redis.CalculateTotalScore(keys[0], value)
				redis.AddToRanking(keys[0], value, score)
				redis.AddVideoToGlobalRanking(fmt.Sprintf("[user-%v] %v", keys[0], value), score)
			} else {
				log.Printf("Error saving data to Redis: %v %v\n", key, value)
			}

			if err != nil {
				log.Printf("Error saving data to Redis: %v\n", err)
			} else {
				fmt.Printf("Successfully saved to Redis: key = %s, value = %s\n", key, value)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v\n", err)
		}
	}
}
