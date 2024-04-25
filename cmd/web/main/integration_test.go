// integration_test.go

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_SearchHospitals(t *testing.T) {
	// Create a new HTTP request to search hospitals in a specific city
	city := "Atyrau"
	searchQuery := "Городская"
	req, err := http.NewRequest("GET", "/search/hospitals/"+city+"?name="+searchQuery, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Make the request to your application
	handler := http.HandlerFunc(searchHospitalsHandler) // Assuming searchHospitalsHandler is in the same package
	handler.ServeHTTP(rr, req)

	// Check the status code is what you expect
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the JSON response
	var hospitalsData PhonebookData
	err = json.Unmarshal(rr.Body.Bytes(), &hospitalsData)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the hospitals were correctly filtered based on the search query
	assert.NotEmpty(t, hospitalsData.Hospitals, "No hospitals found matching the search query")

	// Check if the hospitals contain the expected search query in their names
	for _, hospital := range hospitalsData.Hospitals {
		assert.True(t, strings.Contains(strings.ToLower(hospital.Name), strings.ToLower(searchQuery)), "Hospital name doesn't match search query")
	}
}
