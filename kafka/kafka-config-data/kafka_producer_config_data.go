package kafkaconfigdata

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducerConfigData struct {
	CompressionType   string
	Acks              string
	BatchSize         int
	LingerMs          int
	RequestTimeoutMs  int
	RetryCount        int
	DeliveryTimeoutMs int
}

func NewKafkaProducerConfigData() *KafkaProducerConfigData {
	return &KafkaProducerConfigData{
		CompressionType:   getEnv("KAFKA_PRODUCER_COMPRESSION_TYPE", "snappy"),
		Acks:              getEnv("KAFKA_PRODUCER_ACKS", "all"),
		BatchSize:         getEnvAsInt("KAFKA_PRODUCER_BATCH_SIZE", 16384),
		LingerMs:          getEnvAsInt("KAFKA_PRODUCER_LINGER_MS", 5),
		RequestTimeoutMs:  getEnvAsInt("KAFKA_PRODUCER_REQUEST_TIMEOUT_MS", 30000),
		RetryCount:        getEnvAsInt("KAFKA_PRODUCER_RETRY_COUNT", 3),
		DeliveryTimeoutMs: getEnvAsInt("KAFKA_PRODUCER_DELIVERY_TIMEOUT_MS", 120000),
	}
}

func (k *KafkaProducerConfigData) ToConfigMap(baseConfig kafka.ConfigMap) kafka.ConfigMap {
	config := kafka.ConfigMap{}
	for key, value := range baseConfig {
		config[key] = value
	}

	config["compression.type"] = k.CompressionType
	config["acks"] = k.Acks
	config["batch.size"] = k.BatchSize
	config["linger.ms"] = k.LingerMs
	config["request.timeout.ms"] = k.RequestTimeoutMs
	config["retries"] = k.RetryCount
	config["delivery.timeout.ms"] = k.DeliveryTimeoutMs

	return config
}


