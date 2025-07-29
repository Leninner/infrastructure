package kafkamodel

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   []kafka.Header
	Partition int32
	Offset    kafka.Offset
	Timestamp time.Time
}

func NewMessage(topic string, key []byte, value []byte) *Message {
	return &Message{
		Topic:     topic,
		Key:       key,
		Value:     value,
		Headers:   []kafka.Header{},
		Timestamp: time.Now(),
	}
}

func (m *Message) AddHeader(key, value string) {
	m.Headers = append(m.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (m *Message) SetProtobufValue(msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	m.Value = data
	return nil
}

func (m *Message) GetProtobufValue(msg proto.Message) error {
	return proto.Unmarshal(m.Value, msg)
}
