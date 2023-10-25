package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Проверяем, что функция interpretQueryWithOpenAI коректно находит город среди пользовательского запроса.
func TestInterpretQueryWithOpenAI(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}
		if req.Header.Get("Authorization") != "Bearer "+OPENAI_API_KEY {
			t.Errorf("Expected Authorization Bearer %s, got %s", OPENAI_API_KEY, req.Header.Get("Authorization"))
		}

		response := OpenAIResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{
						Content: "London",
					},
				},
			},
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(response)
	}))
	defer server.Close()

	query := "What's the weather like in London?"
	expectedCity := "London"

	oldDefaultClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	defer func() { http.DefaultClient = oldDefaultClient }()

	city, err := interpretQueryWithOpenAI(query)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if city != expectedCity {
		t.Errorf("Expected city '%s', got '%s'", expectedCity, city)
	}
}
