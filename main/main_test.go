package main

import (
	"testing"
)

func TestFormulateResponse(t *testing.T) {
	testWeatherData := WeatherResponse{
		Main: struct {
			Temp float64 `json:"temp"`
		}{
			Temp: 290.15, // ~17°C
		},
		Weather: []struct {
			Description string `json:"description"`
		}{
			{
				Description: "clear",
			},
		},
	}

	got := formulateResponse("Paris", "What's the temperature in Paris right now?", testWeatherData)
	expected := "The temperature in Paris is 17.00°C."

	if got != expected {
		t.Errorf("got %s; expected %s", got, expected)
	}
}
