package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

var config = sarama.NewConfig()

const (
	LIKE_TOPIC     = "video-likes"
	COMMENT_TOPIC  = "video-comments"
	SHARE_TOPIC    = "video-shares"
	TIMELINE_TOPIC = "video-timeline"
)

func SendMessage(str string, key string, topic string) (bool, error) {
	// Khởi tạo cấu hình Kafka producer
	config.Producer.Return.Successes = true

	// Kết nối tới Kafka broker
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_URL")}, config)
	if err != nil {
		log.Fatal("Error creating producer: ", err)
		return false, err
	}
	defer producer.Close()

	// Gửi tin nhắn
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(str),
		Key:   sarama.StringEncoder(key),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Fatal("Error sending message: ", err)
		return false, err
	}

	return true, nil
}
