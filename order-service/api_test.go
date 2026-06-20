package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProducer struct {
	Topic string
	Key   []byte
	Value []byte
	Err   error
}

func (m *MockProducer) Produce(ctx context.Context, topic string, key []byte, value []byte) error {
	m.Topic = topic
	m.Key = key
	m.Value = value
	return m.Err
}

func TestHealthzEndpoint(t *testing.T) {
	app := &App{Producer: &MockProducer{}}
	
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.healthzHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestOrderEndpoint_Success(t *testing.T) {
	mockProducer := &MockProducer{}
	app := &App{Producer: mockProducer}

	payload := map[string]interface{}{
		"order_id": "123e4567-e89b-12d3-a456-426614174000",
		"item":     "Widget",
		"qty":      5,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.orderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if response["status"] != "Order Transmitted" || response["order_id"] != payload["order_id"] {
		t.Errorf("handler returned unexpected body: %v", response)
	}

	// Verify Kafka producer was called correctly
	if mockProducer.Topic != "orders" {
		t.Errorf("expected topic 'orders', got '%s'", mockProducer.Topic)
	}
	if string(mockProducer.Key) != payload["order_id"] {
		t.Errorf("expected key '%s', got '%s'", payload["order_id"], string(mockProducer.Key))
	}

	var event OrderEvent
	err = json.Unmarshal(mockProducer.Value, &event)
	if err != nil {
		t.Fatalf("failed to parse produced event body: %v", err)
	}

	if event.OrderID != payload["order_id"] || event.Item != payload["item"] || event.Qty != payload["qty"].(int) {
		t.Errorf("produced event data does not match payload: %+v", event)
	}
	if event.Timestamp == "" {
		t.Errorf("expected timestamp to be set in event")
	}
}

func TestOrderEndpoint_InvalidPayload(t *testing.T) {
	app := &App{Producer: &MockProducer{}}

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"Missing order_id", map[string]interface{}{"item": "Widget", "qty": 1}},
		{"Missing item", map[string]interface{}{"order_id": "123", "qty": 1}},
		{"Missing qty", map[string]interface{}{"order_id": "123", "item": "Widget"}},
		{"Invalid qty (zero)", map[string]interface{}{"order_id": "123", "item": "Widget", "qty": 0}},
		{"Invalid qty (negative)", map[string]interface{}{"order_id": "123", "item": "Widget", "qty": -5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.orderHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})
	}
}

func TestOrderEndpoint_ProducerError(t *testing.T) {
	importErr := bytes.NewBufferString("some error") // using an arbitrary error-like condition for mock
	_ = importErr // just placeholder since we can use simple mock error
	
	// Create mock that returns an error
	mockProducer := &MockProducer{Err: http.ErrServerClosed} // any non-nil error works
	app := &App{Producer: mockProducer}

	payload := map[string]interface{}{
		"order_id": "123",
		"item":     "Widget",
		"qty":      5,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.orderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected 500 status on producer error, got %v", status)
	}
}

func TestOrderEndpoint_MethodNotAllowed(t *testing.T) {
	app := &App{Producer: &MockProducer{}}
	req, _ := http.NewRequest("GET", "/order", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.orderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %v", status)
	}
}

func TestSendError(t *testing.T) {
	rr := httptest.NewRecorder()
	sendError(rr, "test error", http.StatusTeapot)
	if rr.Code != http.StatusTeapot {
		t.Errorf("expected 418 Teapot, got %v", rr.Code)
	}
}

func TestOrderEndpoint_InvalidJSON(t *testing.T) {
	app := &App{Producer: &MockProducer{}}
	req, _ := http.NewRequest("POST", "/order", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.orderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid json, got %v", status)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	app := &App{Producer: &MockProducer{}}
	
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	// We will implement app.metricsHandler using promhttp.Handler()
	handler := http.HandlerFunc(app.metricsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("erp_orders_total")) {
		t.Errorf("expected metrics response to contain 'erp_orders_total'")
	}
}
