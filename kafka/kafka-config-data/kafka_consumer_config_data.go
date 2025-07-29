package kafkaconfigdata

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaConsumerConfigData struct {
	GroupID                string
	AutoOffsetReset        string
	SessionTimeoutMs       int
	HeartbeatIntervalMs    int
	MaxPollIntervalMs      int
	MaxPollRecords         int
	MaxPartitionFetchBytes int
	EnableAutoCommit       bool
	AutoCommitIntervalMs   int
	FetchMinBytes          int
	FetchMaxWaitMs         int
	IsolationLevel         string
}

func (k *KafkaConsumerConfigData) ToConfigMap(baseConfig kafka.ConfigMap) kafka.ConfigMap {
	config := baseConfig

	if k.GroupID != "" {
		config["group.id"] = k.GroupID
	}

	if k.AutoOffsetReset != "" {
		config["auto.offset.reset"] = k.AutoOffsetReset
	}

	if k.SessionTimeoutMs > 0 {
		config["session.timeout.ms"] = k.SessionTimeoutMs
	}

	if k.HeartbeatIntervalMs > 0 {
		config["heartbeat.interval.ms"] = k.HeartbeatIntervalMs
	}

	if k.MaxPollIntervalMs > 0 {
		config["max.poll.interval.ms"] = k.MaxPollIntervalMs
	}

	if k.MaxPollRecords > 0 {
		config["max.poll.records"] = k.MaxPollRecords
	}

	if k.MaxPartitionFetchBytes > 0 {
		config["max.partition.fetch.bytes"] = k.MaxPartitionFetchBytes
	}

	config["enable.auto.commit"] = k.EnableAutoCommit

	if k.AutoCommitIntervalMs > 0 {
		config["auto.commit.interval.ms"] = k.AutoCommitIntervalMs
	}

	if k.FetchMinBytes > 0 {
		config["fetch.min.bytes"] = k.FetchMinBytes
	}

	if k.FetchMaxWaitMs > 0 {
		config["fetch.max.wait.ms"] = k.FetchMaxWaitMs
	}

	if k.IsolationLevel != "" {
		config["isolation.level"] = k.IsolationLevel
	}

	return config
}
