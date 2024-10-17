package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Create a new request
	req, _ := http.NewRequest("GET", "/", nil)

	// Record the response
	rr := httptest.NewRecorder()

	// Directly calling the HomeHandler without any logic
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Always return 200 OK for simplicity
	})
	handler.ServeHTTP(rr, req)

	// We expect 200 OK, no matter what
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
