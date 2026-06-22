package proxy

import (
	"errors"
	"testing"
)

func TestRouter_SelectBest(t *testing.T) {
	router := NewRouter()

	// Add endpoints with their AZs
	router.AddEndpoint("broker1:9092", "az-1")
	router.AddEndpoint("broker2:9092", "az-2")
	router.AddEndpoint("broker3:9092", "az-3")

	// Set health statuses
	router.UpdateHealth("broker1:9092", false) // Dead
	router.UpdateHealth("broker2:9092", true)  // Alive
	router.UpdateHealth("broker3:9092", true)  // Alive

	// 1. Test routing when client is in az-1
	// Even though broker1 is in az-1, it's dead. So it should pick broker2 or broker3.
	ep, err := router.SelectBest("az-1")
	if err != nil {
		t.Fatalf("Failed to select best endpoint: %v", err)
	}
	if ep != "broker2:9092" && ep != "broker3:9092" {
		t.Errorf("Expected broker2 or broker3 to be selected, got %s", ep)
	}

	// 2. Test routing when client is in az-2
	// Should pick broker2 since it's alive and in the same AZ
	ep, err = router.SelectBest("az-2")
	if err != nil {
		t.Fatalf("Failed to select best endpoint: %v", err)
	}
	if ep != "broker2:9092" {
		t.Errorf("Expected broker2 to be selected for az-2, got %s", ep)
	}

	// 3. Test when all are dead
	router.UpdateHealth("broker2:9092", false)
	router.UpdateHealth("broker3:9092", false)
	_, err = router.SelectBest("az-1")
	if err == nil || !errors.Is(err, ErrNoHealthyEndpoints) {
		t.Errorf("Expected ErrNoHealthyEndpoints, got %v", err)
	}
}
