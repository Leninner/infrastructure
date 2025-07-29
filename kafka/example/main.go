package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafkaconfig "github.com/leninner/infrastructure/kafka/kafka-config-data"
	kafkaconsumer "github.com/leninner/infrastructure/kafka/kafka-consumer"
	kafkamodel "github.com/leninner/infrastructure/kafka/kafka-model"
	generated "github.com/leninner/infrastructure/kafka/kafka-model/generated"
	kafkaproducer "github.com/leninner/infrastructure/kafka/kafka-producer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	producer *kafkaproducer.Producer
	consumer *kafkaconsumer.Consumer
}

func NewOrderService() (*OrderService, error) {
	baseConfig := &kafkaconfig.KafkaConfigData{
		BootstrapServers:  getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		NumOfPartitions:   3,
		ReplicationFactor: 1,
	}

	producerConfig := &kafkaconfig.KafkaProducerConfigData{
		CompressionType:   getEnv("KAFKA_PRODUCER_COMPRESSION_TYPE", "snappy"),
		Acks:              getEnv("KAFKA_PRODUCER_ACKS", "all"),
		BatchSize:         16384,
		LingerMs:          5,
		RequestTimeoutMs:  30000,
		RetryCount:        3,
		DeliveryTimeoutMs: 120000,
	}

	consumerConfig := &kafkaconfig.KafkaConsumerConfigData{
		GroupID:             getEnv("KAFKA_CONSUMER_GROUP_ID", "order-service-group"),
		AutoOffsetReset:     getEnv("KAFKA_CONSUMER_AUTO_OFFSET_RESET", "earliest"),
		SessionTimeoutMs:    30000,
		HeartbeatIntervalMs: 3000,
		MaxPollIntervalMs:   300000,
		MaxPollRecords:      500,
		EnableAutoCommit:    false,
		FetchMinBytes:       1,
		FetchMaxWaitMs:      500,
	}

	producer, err := kafkaproducer.NewProducer(baseConfig, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	consumer, err := kafkaconsumer.NewConsumer(baseConfig, consumerConfig)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &OrderService{
		producer: producer,
		consumer: consumer,
	}, nil
}

func (s *OrderService) Start() error {
	err := s.consumer.Subscribe([]string{"payment-response", "restaurant-approval-response"})
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	handler := s.createMessageHandler()
	s.consumer.StartConsuming(handler)

	log.Println("Order service started successfully")
	return nil
}

func (s *OrderService) createMessageHandler() kafkaconsumer.MessageHandler {
	return func(msg *kafkamodel.Message) error {
		log.Printf("Received message from topic: %s", msg.Topic)

		switch msg.Topic {
		case "payment-response":
			return s.handlePaymentResponse(msg)
		case "restaurant-approval-response":
			return s.handleRestaurantApprovalResponse(msg)
		default:
			log.Printf("Unknown topic: %s", msg.Topic)
			return nil
		}
	}
}

func (s *OrderService) handlePaymentResponse(msg *kafkamodel.Message) error {
	paymentResponse := &generated.PaymentResponseAvroModel{}
	deserializer := kafkamodel.NewDeserializer()

	err := deserializer.Deserialize(msg.Value, paymentResponse)
	if err != nil {
		return fmt.Errorf("failed to deserialize payment response: %w", err)
	}

	log.Printf("Payment Response - ID: %s, Status: %s, Order: %s",
		paymentResponse.Id, paymentResponse.PaymentStatus, paymentResponse.OrderId)

	if paymentResponse.PaymentStatus == generated.PaymentStatus_PAYMENT_COMPLETED {
		return s.sendRestaurantApprovalRequest(paymentResponse)
	}

	return nil
}

func (s *OrderService) handleRestaurantApprovalResponse(msg *kafkamodel.Message) error {
	approvalResponse := &generated.RestaurantApprovalResponseAvroModel{}
	deserializer := kafkamodel.NewDeserializer()

	err := deserializer.Deserialize(msg.Value, approvalResponse)
	if err != nil {
		return fmt.Errorf("failed to deserialize restaurant approval response: %w", err)
	}

	log.Printf("Restaurant Approval Response - ID: %s, Status: %s, Order: %s",
		approvalResponse.Id, approvalResponse.OrderApprovalStatus, approvalResponse.OrderId)

	return nil
}

func (s *OrderService) sendRestaurantApprovalRequest(paymentResponse *generated.PaymentResponseAvroModel) error {
	approvalRequest := &generated.RestaurantApprovalRequestAvroModel{
		Id:                    generateUUID(),
		SagaId:                paymentResponse.SagaId,
		RestaurantId:          "restaurant-123",
		OrderId:               paymentResponse.OrderId,
		RestaurantOrderStatus: generated.RestaurantOrderStatus_PAID,
		Products: []*generated.Product{
			{Id: "product-1", Quantity: 2},
			{Id: "product-2", Quantity: 1},
		},
		Price:     paymentResponse.Price,
		CreatedAt: timestamppb.Now().Seconds,
	}

	serializer := kafkamodel.NewSerializer()
	data, err := serializer.Serialize(approvalRequest)
	if err != nil {
		return fmt.Errorf("failed to serialize restaurant approval request: %w", err)
	}

	message := kafkamodel.NewMessage("restaurant-approval-request", []byte(approvalRequest.Id), data)
	message.AddHeader("content-type", "application/protobuf")
	message.AddHeader("version", "1.0")

	err = s.producer.ProduceMessage(message)
	if err != nil {
		return fmt.Errorf("failed to produce restaurant approval request: %w", err)
	}

	log.Printf("Sent restaurant approval request for order: %s", approvalRequest.OrderId)
	return nil
}

func (s *OrderService) CreateOrder(customerID, restaurantID string, products []*generated.Product, price string) error {
	paymentRequest := &generated.PaymentRequestAvroModel{
		Id:                 generateUUID(),
		SagaId:             generateUUID(),
		CustomerId:         customerID,
		OrderId:            generateUUID(),
		Price:              []byte(price),
		CreatedAt:          timestamppb.Now().Seconds,
		PaymentOrderStatus: generated.PaymentOrderStatus_PENDING,
	}

	serializer := kafkamodel.NewSerializer()
	data, err := serializer.Serialize(paymentRequest)
	if err != nil {
		return fmt.Errorf("failed to serialize payment request: %w", err)
	}

	message := kafkamodel.NewMessage("payment-request", []byte(paymentRequest.Id), data)
	message.AddHeader("content-type", "application/protobuf")
	message.AddHeader("version", "1.0")

	err = s.producer.ProduceMessage(message)
	if err != nil {
		return fmt.Errorf("failed to produce payment request: %w", err)
	}

	log.Printf("Created order: %s for customer: %s", paymentRequest.OrderId, customerID)
	return nil
}

func (s *OrderService) Stop() {
	log.Println("Stopping order service...")
	s.consumer.Close()
	s.producer.Close()
	log.Println("Order service stopped")
}

func main() {
	service, err := NewOrderService()
	if err != nil {
		log.Fatalf("Failed to create order service: %v", err)
	}
	defer service.Stop()

	err = service.Start()
	if err != nil {
		log.Fatalf("Failed to start order service: %v", err)
	}

	products := []*generated.Product{
		{Id: "pizza-margherita", Quantity: 1},
		{Id: "coke", Quantity: 2},
	}

	err = service.CreateOrder("customer-123", "restaurant-456", products, "25.99")
	if err != nil {
		log.Printf("Failed to create order: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateUUID() string {
	return fmt.Sprintf("uuid-%d", time.Now().UnixNano())
}
