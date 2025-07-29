package kafkaconfigdata

import (
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Environment string

const (
	Development Environment = "dev"
	Staging     Environment = "staging"
	Production  Environment = "prod"
)

type KafkaConfigData struct {
	Environment      Environment
	BootstrapServers string
	SchemaRegistryURL string
	NumOfPartitions   int
	ReplicationFactor int
	SecurityProtocol  string
	SASLMechanism     string
	SASLUsername      string
	SASLPassword      string
}

func NewKafkaConfigData() *KafkaConfigData {
	return &KafkaConfigData{
		Environment:      Environment(getEnv("KAFKA_ENV", "dev")),
		BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		SchemaRegistryURL: getEnv("KAFKA_SCHEMA_REGISTRY_URL", ""),
		NumOfPartitions:   getEnvAsInt("KAFKA_NUM_OF_PARTITIONS", 3),
		ReplicationFactor: getEnvAsInt("KAFKA_REPLICATION_FACTOR", 1),
		SecurityProtocol:  getEnv("KAFKA_SECURITY_PROTOCOL", ""),
		SASLMechanism:     getEnv("KAFKA_SASL_MECHANISM", ""),
		SASLUsername:      getEnv("KAFKA_SASL_USERNAME", ""),
		SASLPassword:      getEnv("KAFKA_SASL_PASSWORD", ""),
	}
}

func (k *KafkaConfigData) ToConfigMap() kafka.ConfigMap {
	config := kafka.ConfigMap{
		"bootstrap.servers": k.BootstrapServers,
	}

	if k.SchemaRegistryURL != "" {
		config["schema.registry.url"] = k.SchemaRegistryURL
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
