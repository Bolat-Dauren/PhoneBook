package models

import (
	"PhoneBook_AP/pkg/drivers"
	"database/sql"
	"testing"
)

func TestCreateUser(t *testing.T) {
	testDB, err := sql.Open("postgres", "user=postgres dbname=finalGo password=0000 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	drivers.DB = testDB

	testUser := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	err = CreateUser(testUser)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query test user: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}
