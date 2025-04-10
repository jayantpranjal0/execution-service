package queue

// import (
// 	"context"
// 	"encoding/json"
// 	"log"

// 	"github.com/segmentio/kafka-go"
// )

// type KafkaClient struct {
// 	writer *kafka.Writer
// 	reader *kafka.Reader
// 	topic  string
// }

// func NewKafkaClient(brokerAddress, topic string) *KafkaClient {
// 	writer := kafka.NewWriter(kafka.WriterConfig{
// 		Brokers: []string{brokerAddress},
// 		Topic:   topic,
// 	})

// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: []string{brokerAddress},
// 		Topic:   topic,
// 		GroupID: "job-execution-group",
// 	})

// 	return &KafkaClient{
// 		writer: writer,
// 		reader: reader,
// 		topic:  topic,
// 	}
// }

// func (kc *KafkaClient) ProduceMessage(ctx context.Context, message interface{}) error {
// 	msg, err := json.Marshal(message)
// 	if err != nil {
// 		return err
// 	}

// 	err = kc.writer.WriteMessages(ctx, kafka.Message{
// 		Value: msg,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (kc *KafkaClient) ConsumeMessages(ctx context.Context) (<-chan []byte, error) {
// 	messages := make(chan []byte)

// 	go func() {
// 		defer close(messages)
// 		for {
// 			msg, err := kc.reader.ReadMessage(ctx)
// 			if err != nil {
// 				log.Println("Error while reading message:", err)
// 				return
// 			}
// 			messages <- msg.Value
// 		}
// 	}()

// 	return messages, nil
// }

// func (kc *KafkaClient) Close() error {
// 	if err := kc.writer.Close(); err != nil {
// 		return err
// 	}
// 	if err := kc.reader.Close(); err != nil {
// 		return err
// 	}
// 	return nil
// }