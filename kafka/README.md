# Kafka Service Migration: Java + Avro → Go + Protocol Buffers

## Overview

This directory contains the complete migration from Java Kafka service using Avro serialization to Go using Protocol Buffers. The migration maintains full compatibility with existing Java services while providing idiomatic Go implementations.

## Architecture

```
infrastructure/kafka/
├── actual/                    # Original Java implementation
│   ├── KafkaConfigData.java
│   ├── KafkaProducerConfigData.java
│   ├── KafkaConsumerConfigData.java
│   └── model/                 # Avro schemas
├── kafka-config-data/         # Go configuration
├── kafka-model/              # Protocol Buffers models
│   ├── resources/proto/      # .proto files
│   └── generated/           # Generated Go code
├── kafka-consumer/          # Consumer implementation
├── kafka-producer/          # Producer implementation
└── README.md
```

## Migration Analysis

### Java → Go Mapping

| Java Component            | Go Equivalent                   | Status     |
| ------------------------- | ------------------------------- | ---------- |
| `KafkaConfigData`         | `kafka_config_data.go`          | ✅ Complete |
| `KafkaProducerConfigData` | `kafka_producer_config_data.go` | ✅ Complete |
| `KafkaConsumerConfigData` | `kafka_consumer_config_data.go` | ✅ Complete |
| Avro Schemas              | Protocol Buffers                | ✅ Complete |
| Spring Boot Config        | Environment Variables           | ✅ Complete |

### Schema Migration

| Avro Schema                           | Protocol Buffer                       | Notes                                            |
| ------------------------------------- | ------------------------------------- | ------------------------------------------------ |
| `CustomerAvroModel`                   | `CustomerAvroModel`                   | UUID → string, added timestamp                   |
| `PaymentRequestAvroModel`             | `PaymentRequestAvroModel`             | Decimal → string, timestamp → protobuf.Timestamp |
| `PaymentResponseAvroModel`            | `PaymentResponseAvroModel`            | Enhanced error handling                          |
| `RestaurantApprovalRequestAvroModel`  | `RestaurantApprovalRequestAvroModel`  | Improved product structure                       |
| `RestaurantApprovalResponseAvroModel` | `RestaurantApprovalResponseAvroModel` | Better status handling                           |

## Usage Examples

### 1. Configuration Setup

```go
package main

import (
    "log"
    kafkaconfig "github.com/leninner/infrastructure/kafka/kafka-config-data"
    kafkaproducer "github.com/leninner/infrastructure/kafka/kafka-producer"
    kafkaconsumer "github.com/leninner/infrastructure/kafka/kafka-consumer"
)

func main() {
    // Base configuration
    baseConfig := &kafkaconfig.KafkaConfigData{
        BootstrapServers:  "localhost:9092",
        NumOfPartitions:   3,
        ReplicationFactor: 1,
    }

    // Producer configuration
    producerConfig := &kafkaconfig.KafkaProducerConfigData{
        CompressionType:     "snappy",
        Acks:               "all",
        BatchSize:          16384,
        LingerMs:           5,
        RequestTimeoutMs:   30000,
        RetryCount:         3,
        DeliveryTimeoutMs:  120000,
    }

    // Consumer configuration
    consumerConfig := &kafkaconfig.KafkaConsumerConfigData{
        GroupID:              "order-service-group",
        AutoOffsetReset:      "earliest",
        SessionTimeoutMs:     30000,
        HeartbeatIntervalMs:  3000,
        MaxPollIntervalMs:    300000,
        MaxPollRecords:       500,
        EnableAutoCommit:     false,
        FetchMinBytes:        1,
        FetchMaxWaitMs:       500,
    }

    // Create producer
    producer, err := kafkaproducer.NewProducer(baseConfig, producerConfig)
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer producer.Close()

    // Create consumer
    consumer, err := kafkaconsumer.NewConsumer(baseConfig, consumerConfig)
    if err != nil {
        log.Fatalf("Failed to create consumer: %v", err)
    }
    defer consumer.Close()
}
```

### 2. Producing Messages

```go
import (
    "github.com/leninner/infrastructure/kafka/kafka-model"
    "github.com/leninner/infrastructure/kafka/kafka-model/generated"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func producePaymentRequest(producer *kafkaproducer.Producer) error {
    // Create protobuf message
    paymentRequest := &generated.PaymentRequestAvroModel{
        Id:                  "123e4567-e89b-12d3-a456-426614174000",
        SagaId:             "456e7890-e89b-12d3-a456-426614174001",
        CustomerId:         "789e0123-e89b-12d3-a456-426614174002",
        OrderId:            "012e3456-e89b-12d3-a456-426614174003",
        Price:              "29.99",
        CreatedAt:          timestamppb.Now(),
        PaymentOrderStatus: generated.PaymentOrderStatus_PENDING,
    }

    // Serialize message
    serializer := kafkamodel.NewSerializer()
    data, err := serializer.Serialize(paymentRequest)
    if err != nil {
        return err
    }

    // Create Kafka message
    message := kafkamodel.NewMessage("payment-request", []byte(paymentRequest.Id), data)
    message.AddHeader("content-type", "application/protobuf")
    message.AddHeader("version", "1.0")

    // Produce message
    return producer.ProduceMessage(message)
}
```

### 3. Consuming Messages

```go
import (
    "context"
    "fmt"
    "github.com/leninner/infrastructure/kafka/kafka-model"
    "github.com/leninner/infrastructure/kafka/kafka-model/generated"
)

func consumePaymentResponse(consumer *kafkaconsumer.Consumer) {
    // Subscribe to topic
    err := consumer.Subscribe([]string{"payment-response"})
    if err != nil {
        log.Fatalf("Failed to subscribe: %v", err)
    }

    // Define message handler
    handler := func(msg *kafkamodel.Message) error {
        // Deserialize message
        paymentResponse := &generated.PaymentResponseAvroModel{}
        deserializer := kafkamodel.NewDeserializer()
        
        err := deserializer.Deserialize(msg.Value, paymentResponse)
        if err != nil {
            return fmt.Errorf("failed to deserialize message: %w", err)
        }

        // Process message
        fmt.Printf("Received payment response: ID=%s, Status=%s\n", 
            paymentResponse.Id, paymentResponse.PaymentStatus)

        return nil
    }

    // Start consuming
    consumer.StartConsuming(handler)

    // Wait for shutdown signal
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Graceful shutdown
    select {
    case <-ctx.Done():
        fmt.Println("Shutting down consumer...")
    }
}
```

## Protocol Buffers Compilation

### Generate Go Code

```bash
# Navigate to proto directory
cd infrastructure/kafka/kafka-model/resources/proto

# Generate Go code
protoc --go_out=../generated \
       --go_opt=paths=source_relative \
       --go-grpc_out=../generated \
       --go-grpc_opt=paths=source_relative \
       *.proto
```

### Required Dependencies

```bash
go get google.golang.org/protobuf/proto
go get google.golang.org/protobuf/types/known/timestamppb
go get github.com/confluentinc/confluent-kafka-go/v2/kafka
```

## Configuration

### Environment Variables

```bash
# Kafka Configuration
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_NUM_OF_PARTITIONS=3
KAFKA_REPLICATION_FACTOR=1

# Producer Configuration
KAFKA_PRODUCER_COMPRESSION_TYPE=snappy
KAFKA_PRODUCER_ACKS=all
KAFKA_PRODUCER_BATCH_SIZE=16384
KAFKA_PRODUCER_LINGER_MS=5
KAFKA_PRODUCER_REQUEST_TIMEOUT_MS=30000
KAFKA_PRODUCER_RETRY_COUNT=3
KAFKA_PRODUCER_DELIVERY_TIMEOUT_MS=120000

# Consumer Configuration
KAFKA_CONSUMER_GROUP_ID=order-service-group
KAFKA_CONSUMER_AUTO_OFFSET_RESET=earliest
KAFKA_CONSUMER_SESSION_TIMEOUT_MS=30000
KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS=3000
KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS=300000
KAFKA_CONSUMER_MAX_POLL_RECORDS=500
KAFKA_CONSUMER_ENABLE_AUTO_COMMIT=false
```

## Key Features

### ✅ Production Ready
- Graceful shutdown with context cancellation
- Comprehensive error handling
- Configurable timeouts and retries
- Batch processing support
- Async message production

### ✅ Protocol Buffers Integration
- Type-safe serialization/deserialization
- Schema evolution support
- Efficient binary format
- Backward compatibility

### ✅ Kafka Best Practices
- Manual offset management
- Proper consumer group handling
- Configurable batch sizes
- Compression support
- Monitoring and observability

### ✅ Go Idioms
- Context-based cancellation
- Error wrapping
- Interface-based design
- Goroutine-safe operations

## Migration Benefits

1. **Performance**: Protocol Buffers are more efficient than Avro
2. **Type Safety**: Compile-time type checking
3. **Schema Evolution**: Better backward compatibility
4. **Go Ecosystem**: Native Go tooling and libraries
5. **Memory Usage**: Lower memory footprint
6. **Serialization Speed**: Faster serialization/deserialization

## Testing

### Unit Tests

```go
func TestPaymentRequestSerialization(t *testing.T) {
    paymentRequest := &generated.PaymentRequestAvroModel{
        Id:     "test-id",
        Price:  "10.99",
        CreatedAt: timestamppb.Now(),
    }

    serializer := kafkamodel.NewSerializer()
    data, err := serializer.Serialize(paymentRequest)
    assert.NoError(t, err)

    deserializer := kafkamodel.NewDeserializer()
    result := &generated.PaymentRequestAvroModel{}
    err = deserializer.Deserialize(data, result)
    assert.NoError(t, err)
    assert.Equal(t, paymentRequest.Id, result.Id)
}
```

### Integration Tests

```go
func TestKafkaIntegration(t *testing.T) {
    // Setup test environment
    // Produce test message
    // Consume and verify
    // Cleanup
}
```

## Monitoring

### Metrics to Track
- Message production rate
- Message consumption rate
- Serialization/deserialization errors
- Consumer lag
- Producer delivery failures
- Connection health

### Health Checks
- Kafka connectivity
- Consumer group status
- Producer delivery status
- Schema registry connectivity

## Troubleshooting

### Common Issues

1. **Serialization Errors**: Verify protobuf schema compatibility
2. **Consumer Lag**: Check batch size and polling intervals
3. **Producer Failures**: Verify acks configuration and timeouts
4. **Connection Issues**: Check bootstrap servers and security settings

### Debug Mode

Enable debug logging by setting environment variable:
```bash
KAFKA_DEBUG=true
```

## Next Steps

1. **Schema Registry**: Integrate with Confluent Schema Registry
2. **Monitoring**: Add Prometheus metrics
3. **Tracing**: Implement distributed tracing
4. **Security**: Add SASL/SSL authentication
5. **Testing**: Comprehensive integration test suite 