# Kafka Module Refactored

## Overview

This document describes the refactored Kafka module that simplifies configuration management and improves maintainability while maintaining Protocol Buffers support.

## Key Improvements

### 1. Simplified Configuration
- **Environment-based configuration**: Automatic configuration loading from environment variables
- **Type-safe configuration**: Strongly typed configuration structures
- **Default values**: Sensible defaults for all configuration options
- **Environment support**: Built-in support for dev, staging, and production environments

### 2. Unified Module Interface
- **Single entry point**: `KafkaModule` provides a unified interface
- **Simplified initialization**: One-line initialization with automatic configuration
- **Protocol Buffers support**: Built-in support for protobuf message serialization
- **Health checks**: Built-in health check functionality

### 3. Improved Maintainability
- **Reduced complexity**: Simplified configuration and usage patterns
- **Better error handling**: Comprehensive error messages and handling
- **Consistent patterns**: Standardized interfaces and patterns
- **Dependency injection**: Maintains DI patterns for testability

## Architecture

```
infrastructure/kafka/
├── kafka-config-data/         # Simplified configuration
│   ├── kafka_config_data.go   # Base Kafka configuration
│   ├── kafka_producer_config_data.go
│   └── kafka_consumer_config_data.go
├── kafka-producer/            # Producer implementation
├── kafka-consumer/            # Consumer implementation
├── kafka-model/              # Protocol Buffers models
├── kafka.module.go           # Unified module interface
└── example/
    └── refactored_usage.go   # Usage examples
```

## Quick Start

### Basic Usage

```go
package main

import (
    kafkamodule "github.com/leninner/infrastructure/kafka"
)

func main() {
    // Initialize with default configuration
    kafkaModule := kafkamodule.NewKafkaModule()
    
    err := kafkaModule.Initialize()
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    defer kafkaModule.Close()

    // Use producer and consumer
    producer := kafkaModule.GetProducer()
    consumer := kafkaModule.GetConsumer()
}
```

### Custom Configuration

```go
package main

import (
    kafkamodule "github.com/leninner/infrastructure/kafka"
    kafkaconfigdata "github.com/leninner/infrastructure/kafka/kafka-config-data"
)

func main() {
    // Custom configuration
    config := &kafkaconfigdata.KafkaConfigData{
        Environment:      kafkaconfigdata.Production,
        BootstrapServers: "kafka-prod:9092",
        SecurityProtocol: "SASL_SSL",
        SASLUsername:     "kafka-user",
        SASLPassword:     "kafka-password",
    }

    producerConfig := &kafkaconfigdata.KafkaProducerConfigData{
        CompressionType: "snappy",
        Acks:           "all",
        BatchSize:      32768,
    }

    consumerConfig := &kafkaconfigdata.KafkaConsumerConfigData{
        GroupID:         "my-service-group",
        AutoOffsetReset: "earliest",
    }

    kafkaModule := kafkamodule.NewKafkaModuleWithConfig(
        config,
        producerConfig,
        consumerConfig,
    )

    err := kafkaModule.Initialize()
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    defer kafkaModule.Close()
}
```

### Protocol Buffers Usage

```go
package main

import (
    kafkamodule "github.com/leninner/infrastructure/kafka"
    "github.com/leninner/infrastructure/kafka/kafka-model/generated"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
    kafkaModule := kafkamodule.NewKafkaModule()
    err := kafkaModule.Initialize()
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    defer kafkaModule.Close()

    // Create protobuf message
    paymentRequest := &generated.PaymentRequestAvroModel{
        Id:     "order-123",
        Price:  "29.99",
        CreatedAt: timestamppb.Now(),
    }

    // Produce protobuf message
    err = kafkaModule.ProduceProtobufMessage(
        "payment-request",
        []byte("order-123"),
        paymentRequest,
    )
    if err != nil {
        log.Fatalf("Failed to produce message: %v", err)
    }

    // Subscribe and consume
    err = kafkaModule.SubscribeToTopics([]string{"payment-request"})
    if err != nil {
        log.Fatalf("Failed to subscribe: %v", err)
    }

    kafkaModule.StartConsuming(func(msg *kafkamodel.Message) error {
        var receivedPaymentRequest generated.PaymentRequestAvroModel
        err := msg.GetProtobufValue(&receivedPaymentRequest)
        if err != nil {
            return err
        }
        
        fmt.Printf("Received: %s\n", receivedPaymentRequest.Id)
        return nil
    })
}
```

## Configuration

### Environment Variables

The module automatically loads configuration from environment variables:

```bash
# Base Kafka Configuration
KAFKA_ENV=dev
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_NUM_OF_PARTITIONS=3
KAFKA_REPLICATION_FACTOR=1

# Security (optional)
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=PLAIN
KAFKA_SASL_USERNAME=kafka-user
KAFKA_SASL_PASSWORD=kafka-password

# Producer Configuration
KAFKA_PRODUCER_COMPRESSION_TYPE=snappy
KAFKA_PRODUCER_ACKS=all
KAFKA_PRODUCER_BATCH_SIZE=16384
KAFKA_PRODUCER_LINGER_MS=5
KAFKA_PRODUCER_REQUEST_TIMEOUT_MS=30000
KAFKA_PRODUCER_RETRY_COUNT=3
KAFKA_PRODUCER_DELIVERY_TIMEOUT_MS=120000

# Consumer Configuration
KAFKA_CONSUMER_GROUP_ID=default-group
KAFKA_CONSUMER_AUTO_OFFSET_RESET=earliest
KAFKA_CONSUMER_SESSION_TIMEOUT_MS=30000
KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS=3000
KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS=300000
KAFKA_CONSUMER_MAX_POLL_RECORDS=500
KAFKA_CONSUMER_FETCH_MIN_BYTES=1
KAFKA_CONSUMER_FETCH_MAX_WAIT_MS=500
```

### Environment-Specific Configurations

```go
// Development
config := &kafkaconfigdata.KafkaConfigData{
    Environment:      kafkaconfigdata.Development,
    BootstrapServers: "localhost:9092",
    NumOfPartitions:  3,
    ReplicationFactor: 1,
}

// Staging
config := &kafkaconfigdata.KafkaConfigData{
    Environment:      kafkaconfigdata.Staging,
    BootstrapServers: "kafka-staging:9092",
    NumOfPartitions:  6,
    ReplicationFactor: 2,
}

// Production
config := &kafkaconfigdata.KafkaConfigData{
    Environment:      kafkaconfigdata.Production,
    BootstrapServers: "kafka-prod:9092",
    NumOfPartitions:  12,
    ReplicationFactor: 3,
    SecurityProtocol: "SASL_SSL",
    SASLMechanism:    "PLAIN",
    SASLUsername:     "kafka-user",
    SASLPassword:     "kafka-password",
}
```

## API Reference

### KafkaModule

```go
type KafkaModule struct {
    // ... internal fields
}

// Constructors
func NewKafkaModule() *KafkaModule
func NewKafkaModuleWithConfig(config, producerConfig, consumerConfig) *KafkaModule

// Initialization
func (m *KafkaModule) Initialize() error
func (m *KafkaModule) InitializeProducer() error
func (m *KafkaModule) InitializeConsumer() error

// Message Production
func (m *KafkaModule) ProduceMessage(topic, key, value, headers) error
func (m *KafkaModule) ProduceProtobufMessage(topic, key, message) error

// Message Consumption
func (m *KafkaModule) SubscribeToTopics(topics) error
func (m *KafkaModule) StartConsuming(handler) error

// Health and Management
func (m *KafkaModule) Health() error
func (m *KafkaModule) Close()

// Accessors
func (m *KafkaModule) GetProducer() *Producer
func (m *KafkaModule) GetConsumer() *Consumer
func (m *KafkaModule) GetConfig() *KafkaConfigData
```

## Migration from Old Module

### Before (Old Module)
```go
// Complex configuration setup
baseConfig := &kafkaconfigdata.KafkaConfigData{
    BootstrapServers: "localhost:9092",
}

producerConfig := &kafkaconfigdata.KafkaProducerConfigData{
    CompressionType: "snappy",
    Acks:           "all",
}

// Manual producer creation
producer, err := kafkaproducer.NewProducer(baseConfig, producerConfig)
if err != nil {
    log.Fatalf("Failed to create producer: %v", err)
}
defer producer.Close()

// Manual consumer creation
consumer, err := kafkaconsumer.NewConsumer(baseConfig, consumerConfig)
if err != nil {
    log.Fatalf("Failed to create consumer: %v", err)
}
defer consumer.Close()
```

### After (Refactored Module)
```go
// Simple initialization
kafkaModule := kafkamodule.NewKafkaModule()
err := kafkaModule.Initialize()
if err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
defer kafkaModule.Close()

// Direct access to producer and consumer
producer := kafkaModule.GetProducer()
consumer := kafkaModule.GetConsumer()
```

## Benefits

### 1. Reduced Complexity
- **50% less code** for basic setup
- **Simplified configuration** with sensible defaults
- **Unified interface** for all Kafka operations

### 2. Improved Developer Experience
- **Type-safe configuration** prevents runtime errors
- **Environment-aware** configuration loading
- **Better error messages** with context

### 3. Production Ready
- **Health checks** for monitoring
- **Graceful shutdown** with proper cleanup
- **Protocol Buffers support** for efficient serialization

### 4. Maintainability
- **Consistent patterns** across the codebase
- **Dependency injection** for testability
- **Clear separation** of concerns

## Testing

```go
func TestKafkaModule(t *testing.T) {
    kafkaModule := kafkamodule.NewKafkaModule()
    
    err := kafkaModule.Initialize()
    assert.NoError(t, err)
    defer kafkaModule.Close()

    // Test health check
    err = kafkaModule.Health()
    assert.NoError(t, err)

    // Test message production
    err = kafkaModule.ProduceMessage(
        "test-topic",
        []byte("test-key"),
        []byte("test-value"),
        nil,
    )
    assert.NoError(t, err)
}
```

## Best Practices

1. **Use environment variables** for configuration in production
2. **Initialize once** and reuse the module instance
3. **Always call Close()** to ensure proper cleanup
4. **Use Protocol Buffers** for message serialization when possible
5. **Implement health checks** in your application
6. **Handle errors gracefully** with proper logging

## Troubleshooting

### Common Issues

1. **Connection failures**: Check bootstrap servers and network connectivity
2. **Authentication errors**: Verify SASL credentials and security protocol
3. **Serialization errors**: Ensure protobuf schemas are compatible
4. **Consumer lag**: Adjust batch size and polling intervals

### Debug Mode

Enable debug logging by setting environment variable:
```bash
KAFKA_DEBUG=true
```

## Next Steps

1. **Schema Registry**: Integrate with Confluent Schema Registry
2. **Monitoring**: Add Prometheus metrics
3. **Tracing**: Implement distributed tracing
4. **Security**: Add more authentication methods
5. **Testing**: Comprehensive integration test suite 