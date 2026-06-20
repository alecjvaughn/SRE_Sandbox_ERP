package main

import (
	"encoding/json"
	"net/http"
)

type OrderRequest struct {
	OrderID string `json:"order_id"`
	Item    string `json:"item"`
	Qty     int    `json:"qty"`
}

type OrderResponse struct {
	Status  string `json:"status,omitempty"`
	OrderID string `json:"order_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
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

	// Mock response for now (Kafka producer will be added later)
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
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/order", orderHandler)
	http.ListenAndServe(":8080", nil)
}
