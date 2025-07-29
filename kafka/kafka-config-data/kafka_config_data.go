package kafkaconfigdata

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaConfigData struct {
	BootstrapServers  string
	SchemaRegistryURL string
	NumOfPartitions   int
	ReplicationFactor int
	SecurityProtocol  string
	SASLMechanism     string
	SASLUsername      string
	SASLPassword      string
}

func (k *KafkaConfigData) ToConfigMap() kafka.ConfigMap {
	config := kafka.ConfigMap{
		"bootstrap.servers": k.BootstrapServers,
	}

	if k.SecurityProtocol != "" {
		config["security.protocol"] = k.SecurityProtocol
	}

	if k.SASLMechanism != "" {
		config["sasl.mechanism"] = k.SASLMechanism
	}

	if k.SASLUsername != "" {
		config["sasl.username"] = k.SASLUsername
	}

	if k.SASLPassword != "" {
		config["sasl.password"] = k.SASLPassword
	}

	return config
}
