package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndToEndInvalidRegistrationPage(t *testing.T) {
	registrationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/register" {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "page not found", http.StatusNotFound)
	})

	server := httptest.NewServer(registrationHandler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/invalid_registration")
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code, got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
