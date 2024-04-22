package models

import (
	"PhoneBook_AP/pkg/drivers"
	"database/sql"
	"testing"
)

func TestIntegrationGetHashedPassword(t *testing.T) {
	testDB, err := sql.Open("postgres", "user=postgres dbname=finalGo password=0000 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	drivers.DB = testDB

	_, err = testDB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", "testuser", "test@example.com", "testpassword")
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}
	hashedPassword, err := GetHashedPassword("testuser")
	if err != nil {
		t.Fatalf("failed to get hashed password: %v", err)
	}
	expectedHashedPassword := "$2a$10$"

	if hashedPassword[:7] != expectedHashedPassword {
		t.Errorf("expected hashed password %s, got %s", expectedHashedPassword, hashedPassword)
	}
}
