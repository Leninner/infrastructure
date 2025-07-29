package kafkamodel

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Serializer struct{}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) Serialize(msg proto.Message) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize protobuf message: %w", err)
	}

	return data, nil
}

func (s *Serializer) Deserialize(data []byte, msg proto.Message) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("data is nil or empty")
	}

	if msg == nil {
		return fmt.Errorf("target message is nil")
	}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("failed to deserialize protobuf message: %w", err)
	}

	return nil
}

type Deserializer struct{}

func NewDeserializer() *Deserializer {
	return &Deserializer{}
}

func (d *Deserializer) Deserialize(data []byte, msg proto.Message) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("data is nil or empty")
	}

	if msg == nil {
		return fmt.Errorf("target message is nil")
	}

	err := proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("failed to deserialize protobuf message: %w", err)
	}

	return nil
}
