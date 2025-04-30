package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type exchangeResponse struct {
	Result         string  `json:"result"`
	ConversionRate float64 `json:"conversion_rate"`
}

// FetchRate calls the ExchangeRate-API to get the latest conversion rate between two currencies.
func FetchRate(baseURL, apiKey, from, to string) (float64, error) {
	url := fmt.Sprintf("%s/%s/pair/%s/%s", baseURL, apiKey, from, to)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	var result exchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if result.Result != "success" {
		return 0, fmt.Errorf("API error: result=%s", result.Result)
	}

	return result.ConversionRate, nil
}
