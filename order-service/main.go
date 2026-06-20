package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type OrderRequest struct {
	OrderID string `json:"order_id"`
	Item    string `json:"item"`
	Qty     int    `json:"qty"`
}

type OrderEvent struct {
	OrderID   string `json:"order_id"`
	Item      string `json:"item"`
	Qty       int    `json:"qty"`
	Timestamp string `json:"timestamp"`
}

type OrderResponse struct {
	Status  string `json:"status,omitempty"`
	OrderID string `json:"order_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

// EventProducer interface allows mocking the Kafka producer
type EventProducer interface {
	Produce(ctx context.Context, topic string, key []byte, value []byte) error
}

// KafkaProducer wraps the kgo.Client
type KafkaProducer struct {
	Client *kgo.Client
}

func (k *KafkaProducer) Produce(ctx context.Context, topic string, key []byte, value []byte) error {
	record := &kgo.Record{Topic: topic, Key: key, Value: value}
	if err := k.Client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return err
	}
	return nil
}

type App struct {
	Producer EventProducer
}

func (app *App) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app *App) orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		sendError(w, "Missing order_id", http.StatusBadRequest)
		return
	}
	if req.Item == "" {
		sendError(w, "Missing item", http.StatusBadRequest)
		return
	}
	if req.Qty <= 0 {
		sendError(w, "Invalid qty", http.StatusBadRequest)
		return
	}

	event := OrderEvent{
		OrderID:   req.OrderID,
		Item:      req.Item,
		Qty:       req.Qty,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	eventBytes, _ := json.Marshal(event)

	err := app.Producer.Produce(r.Context(), "orders", []byte(req.OrderID), eventBytes)
	if err != nil {
		log.Printf("Failed to produce event: %v", err)
		sendError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(OrderResponse{
		Status:  "Order Transmitted",
		OrderID: req.OrderID,
	})
}

func sendError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(OrderResponse{Error: msg})
}

func main() {
	brokers := []string{"kafka:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		brokers = []string{envBrokers}
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		log.Fatalf("unable to create kafka client: %v", err)
	}
	defer client.Close()

	app := &App{
		Producer: &KafkaProducer{Client: client},
	}

	http.HandleFunc("/healthz", app.healthzHandler)
	http.HandleFunc("/order", app.orderHandler)
	
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
