package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	OPENWEATHER_API_KEY = "e466bef46d86e7b3c52399fa85c0e758"
	OPENAI_API_KEY      = "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"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func userInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ask a weather-related question:")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	city, err := interpretQueryWithOpenAI(query)
	if err != nil {
		fmt.Println("Error interpreting query:", err)
		return
	}

	weatherData, err := getWeather(city)
	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return
	}

	response, err := formulateResponseWithOpenAI(city, query, weatherData)
	if err != nil {
		fmt.Println("Error formulating response:", err)
		return
	}
	fmt.Println(response)
}

func formulateResponseWithOpenAI(city, query string, weatherData WeatherResponse) (string, error) {
	tempCelsius := weatherData.Main.Temp - 273.15
	weatherDescription := "Unknown"
	if len(weatherData.Weather) > 0 {
		weatherDescription = weatherData.Weather[0].Description
	}

	systemMessage := "You are a helpful and friendly assistant. Given the current weather data, please provide an informative and conversational response to the user's question."
	userMessage := fmt.Sprintf("The temperature in %s is %.2f°C, and the weather is described as %s. Question: %s", city, tempCelsius, weatherDescription, query)

	messages := []Message{
		{Role: "system", Content: systemMessage},
		{Role: "user", Content: userMessage},
	}

	payload := map[string]interface{}{
		"model":    "gpt-4",
		"messages": messages,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	var response OpenAIResponse

	json.Unmarshal(bodyResp, &response)

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no content from openAI")
}

func interpretQueryWithOpenAI(query string) (string, error) {

	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant. Please provide only the city name when asked about weather. Don't say anything else"},
		{Role: "user", Content: query},
	}

	payload := map[string]interface{}{
		"model":    "gpt-4",
		"messages": messages,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	var response OpenAIResponse

	fmt.Println("OpenAI Interpret Response:", string(bodyResp))

	json.Unmarshal(bodyResp, &response)

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("OpenAI API did not return expected choices")
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}
	return response.Choices[0].Message.Content, nil
}

func getWeather(city string) (WeatherResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, OPENWEATHER_API_KEY))
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var weatherData WeatherResponse
	err = json.Unmarshal(data, &weatherData)
	if err != nil {
		return WeatherResponse{}, err
	}
	return weatherData, nil
}

func main() {
	userInput()
}
