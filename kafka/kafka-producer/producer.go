package kafkaproducer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafkaconfigdata "github.com/leninner/infrastructure/kafka/kafka-config-data"
	kafkamodel "github.com/leninner/infrastructure/kafka/kafka-model"
)

type Producer struct {
	producer *kafka.Producer
	config   *kafkaconfigdata.KafkaProducerConfigData
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewProducer(baseConfig *kafkaconfigdata.KafkaConfigData, producerConfig *kafkaconfigdata.KafkaProducerConfigData) (*Producer, error) {
	configMap := producerConfig.ToConfigMap(baseConfig.ToConfigMap())

	producer, err := kafka.NewProducer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Producer{
		producer: producer,
		config:   producerConfig,
		ctx:      ctx,
		cancel:   cancel,
	}

	p.wg.Add(1)
	go p.handleEvents()

	return p, nil
}

func (p *Producer) handleEvents() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case e := <-p.producer.Events():
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			case kafka.Error:
				fmt.Printf("Producer error: %v\n", ev)
			}
		}
	}
}

func (p *Producer) Produce(topic string, key []byte, value []byte, headers []kafka.Header) error {
	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
		Headers:        headers,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		msg := e.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
		}
	case <-time.After(time.Duration(p.config.DeliveryTimeoutMs) * time.Millisecond):
		return fmt.Errorf("delivery timeout")
	case <-p.ctx.Done():
		return fmt.Errorf("producer context cancelled")
	}

	return nil
}

func (p *Producer) ProduceMessage(message *kafkamodel.Message) error {
	return p.Produce(message.Topic, message.Key, message.Value, message.Headers)
}

func (p *Producer) ProduceAsync(topic string, key []byte, value []byte, headers []kafka.Header) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
		Headers:        headers,
	}, nil)
}

func (p *Producer) Flush(timeoutMs int) int {
	return p.producer.Flush(timeoutMs)
}

func (p *Producer) Close() {
	p.cancel()
	p.wg.Wait()
	p.producer.Close()
}
