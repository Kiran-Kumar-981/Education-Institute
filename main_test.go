package main

import (
	// Replace with your actual package path
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	testDSN = "UserName:Password@tcp(127.0.0.1:3306)/test_mydata" // Test database DSN
)

// TestDatabase initializes the database connection and checks if it's reachable.
func TestDatabase(t *testing.T) {
	dataBase, err := sql.Open("mysql", testDSN)
	if err != nil {
		t.Fatalf("Error opening database connection: %v", err)
	}
	defer dataBase.Close()

	err = dataBase.Ping()
	if err != nil {
		t.Fatalf("Error pinging database: %v", err)
	}

	t.Logf("Successfully connected to the test database with DSN: %s", testDSN)
}

// TestAdmissionsHandler tests the admissions handler function.
func TestAdmissionsHandler(t *testing.T) {
	// Initialize a Gin router
	router := gin.Default()
	router.LoadHTMLGlob("templete/*.html")

	// Mock database setup
	dataBase, _ := sql.Open("mysql", testDSN)
	defer dataBase.Close()

	// Mock admission form data
	body := strings.NewReader("Name=John&FatherName=Doe&Qualification=Bachelor&Email=johndoe@example.com&PhNumber=1234567890&Course=Engineering&Address=123 Main St&Duration=4 years&Fee=5000&BatchTiming=Morning&FeePaid=2500")

	// Create a mock HTTP request
	req, err := http.NewRequest("POST", "/admissionForm", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	// Check if the inserted.html template is rendered
	if !strings.Contains(w.Body.String(), "inserted.html") {
		t.Errorf("Expected response body to contain 'inserted.html'; got %s", w.Body.String())
	}
}
