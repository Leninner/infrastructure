package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafkaconfigdata "github.com/leninner/infrastructure/kafka/kafka-config-data"
	kafkaconsumer "github.com/leninner/infrastructure/kafka/kafka-consumer"
	kafkamodel "github.com/leninner/infrastructure/kafka/kafka-model"
	kafkaproducer "github.com/leninner/infrastructure/kafka/kafka-producer"
	"google.golang.org/protobuf/proto"
)

type KafkaModule struct {
	config         *kafkaconfigdata.KafkaConfigData
	producerConfig *kafkaconfigdata.KafkaProducerConfigData
	consumerConfig *kafkaconfigdata.KafkaConsumerConfigData
	producer       *kafkaproducer.Producer
	consumer       *kafkaconsumer.Consumer
}

func NewKafkaModule() *KafkaModule {
	return &KafkaModule{
		config:         kafkaconfigdata.NewKafkaConfigData(),
		producerConfig: kafkaconfigdata.NewKafkaProducerConfigData(),
		consumerConfig: kafkaconfigdata.NewKafkaConsumerConfigData(),
	}
}

func NewKafkaModuleWithConfig(
	config *kafkaconfigdata.KafkaConfigData,
	producerConfig *kafkaconfigdata.KafkaProducerConfigData,
	consumerConfig *kafkaconfigdata.KafkaConsumerConfigData,
) *KafkaModule {
	return &KafkaModule{
		config:         config,
		producerConfig: producerConfig,
		consumerConfig: consumerConfig,
	}
}

func (m *KafkaModule) InitializeProducer() error {
	producer, err := kafkaproducer.NewProducer(m.config, m.producerConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize producer: %w", err)
	}
	m.producer = producer
	return nil
}

func (m *KafkaModule) InitializeConsumer() error {
	consumer, err := kafkaconsumer.NewConsumer(m.config, m.consumerConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize consumer: %w", err)
	}
	m.consumer = consumer
	return nil
}

func (m *KafkaModule) Initialize() error {
	if err := m.InitializeProducer(); err != nil {
		return err
	}
	if err := m.InitializeConsumer(); err != nil {
		m.producer.Close()
		return err
	}
	return nil
}

func (m *KafkaModule) GetProducer() *kafkaproducer.Producer {
	return m.producer
}

func (m *KafkaModule) GetConsumer() *kafkaconsumer.Consumer {
	return m.consumer
}

func (m *KafkaModule) ProduceProtobufMessage(topic string, key []byte, message proto.Message) error {
	if m.producer == nil {
		return fmt.Errorf("producer not initialized")
	}

	kafkaMessage := kafkamodel.NewMessage(topic, key, nil)
	
	if err := kafkaMessage.SetProtobufValue(message); err != nil {
		return fmt.Errorf("failed to set protobuf value: %w", err)
	}

	return m.producer.ProduceMessage(kafkaMessage)
}

func (m *KafkaModule) ProduceMessage(topic string, key []byte, value []byte, headers []kafka.Header) error {
	if m.producer == nil {
		return fmt.Errorf("producer not initialized")
	}

	return m.producer.Produce(topic, key, value, headers)
}

func (m *KafkaModule) SubscribeToTopics(topics []string) error {
	if m.consumer == nil {
		return fmt.Errorf("consumer not initialized")
	}

	return m.consumer.Subscribe(topics)
}

func (m *KafkaModule) StartConsuming(handler kafkaconsumer.MessageHandler) {
	if m.consumer == nil {
		return
	}

	m.consumer.StartConsuming(handler)
}

func (m *KafkaModule) Health() error {
	if m.producer == nil {
		return fmt.Errorf("producer not initialized")
	}

	if m.consumer == nil {
		return fmt.Errorf("consumer not initialized")
	}

	return nil
}

func (m *KafkaModule) Close() {
	if m.producer != nil {
		m.producer.Close()
	}
	if m.consumer != nil {
		m.consumer.Close()
	}
}

func (m *KafkaModule) GetConfig() *kafkaconfigdata.KafkaConfigData {
	return m.config
}

func (m *KafkaModule) GetProducerConfig() *kafkaconfigdata.KafkaProducerConfigData {
	return m.producerConfig
}

func (m *KafkaModule) GetConsumerConfig() *kafkaconfigdata.KafkaConsumerConfigData {
	return m.consumerConfig
} 