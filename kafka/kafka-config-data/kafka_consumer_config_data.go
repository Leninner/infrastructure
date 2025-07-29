package kafkaconfigdata

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumerConfigData struct {
	GroupID             string
	AutoOffsetReset     string
	SessionTimeoutMs    int
	HeartbeatIntervalMs int
	MaxPollIntervalMs   int
	MaxPollRecords      int
	EnableAutoCommit    bool
	FetchMinBytes       int
	FetchMaxWaitMs      int
}

func NewKafkaConsumerConfigData() *KafkaConsumerConfigData {
	return &KafkaConsumerConfigData{
		GroupID:             getEnv("KAFKA_CONSUMER_GROUP_ID", "default-group"),
		AutoOffsetReset:     getEnv("KAFKA_CONSUMER_AUTO_OFFSET_RESET", "earliest"),
		SessionTimeoutMs:    getEnvAsInt("KAFKA_CONSUMER_SESSION_TIMEOUT_MS", 30000),
		HeartbeatIntervalMs: getEnvAsInt("KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS", 3000),
		MaxPollIntervalMs:   getEnvAsInt("KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS", 300000),
		MaxPollRecords:      getEnvAsInt("KAFKA_CONSUMER_MAX_POLL_RECORDS", 500),
		EnableAutoCommit:    false,
		FetchMinBytes:       getEnvAsInt("KAFKA_CONSUMER_FETCH_MIN_BYTES", 1),
		FetchMaxWaitMs:      getEnvAsInt("KAFKA_CONSUMER_FETCH_MAX_WAIT_MS", 500),
	}
}

func (k *KafkaConsumerConfigData) ToConfigMap(baseConfig kafka.ConfigMap) kafka.ConfigMap {
	config := kafka.ConfigMap{}
	for key, value := range baseConfig {
		config[key] = value
	}

	config["group.id"] = k.GroupID
	config["auto.offset.reset"] = k.AutoOffsetReset
	config["session.timeout.ms"] = k.SessionTimeoutMs
	config["heartbeat.interval.ms"] = k.HeartbeatIntervalMs
	config["max.poll.interval.ms"] = k.MaxPollIntervalMs
	config["max.poll.records"] = k.MaxPollRecords
	config["enable.auto.commit"] = k.EnableAutoCommit
	config["fetch.min.bytes"] = k.FetchMinBytes
	config["fetch.max.wait.ms"] = k.FetchMaxWaitMs

	return config
}
