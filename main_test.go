package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/albums", getAlbums)
	r.GET("/albums/:id", getAlbum)
	r.POST("/create", createAlbums)

	return r
}

func TestGetAlbums(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/albums", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []album
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	if len(response) != 3 {
		t.Fatalf("expected 3 albums, got %d", len(response))
	}
}

func TestGetAlbumByID(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		id           string
		expectedCode int
	}{
		{"1", http.StatusOK},
		{"2", http.StatusOK},
		{"999", http.StatusNotFound},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("GET", "/albums/"+tc.id, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tc.expectedCode {
			t.Errorf("for ID %s expected %d, got %d", tc.id, tc.expectedCode, w.Code)
		}
	}
}

func TestCreateAlbum(t *testing.T) {
	router := setupRouter()

	newAlbum := album{
		ID:     "10",
		Title:  "Test Album",
		Artist: "Tester",
		Price:  9.99,
	}

	body, _ := json.Marshal(newAlbum)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}

	var response album
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("error unmarshalling json: %v", err)
	}

	if response.ID != newAlbum.ID {
		t.Errorf("expected ID %s, got %s", newAlbum.ID, response.ID)
	}
}
