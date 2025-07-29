package kafkaconfigdata

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaProducerConfigData struct {
	CompressionType     string
	Acks                string
	BatchSize           int
	LingerMs            int
	RequestTimeoutMs    int
	RetryCount          int
	EnableIdempotence   bool
	MaxInFlightRequests int
	DeliveryTimeoutMs   int
}

func (k *KafkaProducerConfigData) ToConfigMap(baseConfig kafka.ConfigMap) kafka.ConfigMap {
	config := baseConfig

	if k.CompressionType != "" {
		config["compression.type"] = k.CompressionType
	}

	if k.Acks != "" {
		config["acks"] = k.Acks
	}

	if k.BatchSize > 0 {
		config["batch.size"] = k.BatchSize
	}

	if k.LingerMs > 0 {
		config["linger.ms"] = k.LingerMs
	}

	if k.RequestTimeoutMs > 0 {
		config["request.timeout.ms"] = k.RequestTimeoutMs
	}

	if k.RetryCount > 0 {
		config["retries"] = k.RetryCount
	}

	if k.EnableIdempotence {
		config["enable.idempotence"] = true
	}

	if k.MaxInFlightRequests > 0 {
		config["max.in.flight.requests.per.connection"] = k.MaxInFlightRequests
	}

	if k.DeliveryTimeoutMs > 0 {
		config["delivery.timeout.ms"] = k.DeliveryTimeoutMs
	}

	return config
}
