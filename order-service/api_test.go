package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthzHandler)

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
	handler := http.HandlerFunc(orderHandler)

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
}

func TestOrderEndpoint_InvalidPayload(t *testing.T) {
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
			handler := http.HandlerFunc(orderHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}

			var response map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("failed to parse response body: %v", err)
			}

			if _, ok := response["error"]; !ok {
				t.Errorf("expected error message in response body: %v", response)
			}
		})
	}
}
