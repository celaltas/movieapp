package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"movieexample.com/rating/pkg/model"
)

func main() {
	fmt.Println("Creating a Kafka producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratingsdata.json"
	fmt.Println("Reading rating events from file " + fileName)

	ratingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceRatingEvents(producer, topic, ratingEvents); err != nil {
		panic(err)
	}
	const timeout = 10 * time.Second
	fmt.Println("Waiting " + timeout.String() + " until allevents get produced")
	flushedNum := producer.Flush(int(timeout.Milliseconds()))
	fmt.Println("Flushed " + fmt.Sprint(flushedNum) + " events")

}

func readRatingEvents(filename string) ([]model.RatingEvent, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ratingEvents []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratingEvents); err != nil {
		return nil, err
	}
	return ratingEvents, nil
}

func produceRatingEvents(producer *kafka.Producer, topic string, ratingEvents []model.RatingEvent) error {
	for _, ratingEvent := range ratingEvents {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: encodedEvent,
		}
		if err := producer.Produce(message, nil); err != nil {
			return err
		}
	}
	return nil
}
