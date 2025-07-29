package kafkaconsumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	kafkaconfigdata "github.com/leninner/infrastructure/kafka/kafka-config-data"
	kafkamodel "github.com/leninner/infrastructure/kafka/kafka-model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
	config   *kafkaconfigdata.KafkaConsumerConfigData
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

type MessageHandler func(*kafkamodel.Message) error

func NewConsumer(baseConfig *kafkaconfigdata.KafkaConfigData, consumerConfig *kafkaconfigdata.KafkaConsumerConfigData) (*Consumer, error) {
	configMap := consumerConfig.ToConfigMap(baseConfig.ToConfigMap())

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer: consumer,
		config:   consumerConfig,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (c *Consumer) Subscribe(topics []string) error {
	return c.consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) SubscribePattern(pattern string) error {
	return c.consumer.Subscribe(pattern, nil)
}

func (c *Consumer) Poll(timeoutMs int) (*kafkamodel.Message, error) {
	msg, err := c.consumer.ReadMessage(time.Duration(timeoutMs) * time.Millisecond)
	if err != nil {
		if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return &kafkamodel.Message{
		Topic:     *msg.TopicPartition.Topic,
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   msg.Headers,
		Partition: msg.TopicPartition.Partition,
		Offset:    msg.TopicPartition.Offset,
		Timestamp: msg.Timestamp,
	}, nil
}

func (c *Consumer) StartConsuming(handler MessageHandler) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				msg, err := c.Poll(1000)
				if err != nil {
					fmt.Printf("Error polling message: %v\n", err)
					continue
				}
				if msg == nil {
					continue
				}

				if err := handler(msg); err != nil {
					fmt.Printf("Error handling message: %v\n", err)
				} else {
					if _, err := c.CommitMessage(msg); err != nil {
						fmt.Printf("Error committing message: %v\n", err)
					}
				}
			}
		}
	}()
}

func (c *Consumer) Commit() ([]kafka.TopicPartition, error) {
	return c.consumer.Commit()
}

func (c *Consumer) CommitMessage(message *kafkamodel.Message) ([]kafka.TopicPartition, error) {
	topicPartition := kafka.TopicPartition{
		Topic:     &message.Topic,
		Partition: message.Partition,
		Offset:    message.Offset,
	}

	return c.consumer.CommitMessage(&kafka.Message{
		TopicPartition: topicPartition,
	})
}

func (c *Consumer) Close() error {
	c.cancel()
	c.wg.Wait()
	return c.consumer.Close()
}