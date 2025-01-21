package simetrics

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"none.rks/config/config"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

var tokenCache *Token

func getToken() (*Token, error) {
	authURL := config.AppConfig.Siurl + "/restapi/v1/tenants/" + config.AppConfig.TenantId + "/token"

	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", config.AppConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to retrieve token, status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	result, ok := data["result"].(map[string]interface{})
	if !ok {
		log.Fatal("Error parsing 'result' from JSON")
	}

	token, ok := result["token"].(string)
	if !ok {
		log.Fatal("Error parsing 'token' from JSON")
	}

	expiration, ok := result["expiration"].(float64) // Use float64 for numbers
	if !ok {
		log.Fatal("Error parsing 'expiration' from JSON")
	}

	return &Token{
		Value:     token,
		ExpiresAt: time.Unix(0, int64(expiration)*int64(time.Millisecond)).UTC(),
	}, nil
}

func FetchData() ([]interface{}, error) {
	if tokenCache == nil || time.Now().After(tokenCache.ExpiresAt) {
		var err error
		tokenCache, err = getToken()
		if err != nil {
			return nil, err
		}
	}

	baseURL := config.AppConfig.Siurl + "/restapi/v1/tenants/" + config.AppConfig.TenantId + "/storage-systems/" + config.AppConfig.SystemIds[0] + "/metrics"

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	params := url.Values{}
	for _, metrics := range config.AppConfig.Metrics {
		params.Add("types", metrics)
	}

	params.Add("duration", "1h")

	parsedURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-token", tokenCache.Value)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		tokenCache, err = getToken() // Refresh the token
		if err != nil {
			return nil, err
		}
		return FetchData()
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		log.Fatal("Error parsing 'data' from JSON")
	}

	return data, nil
}
