package neato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func httpGet(uri string, header *url.Values, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if header != nil {
		for k, vv := range *header {
			for _, v := range vv {
				req.Header.Set(k, v)
			}
		}
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("expected HTTP 2xx/3xx, got %s", resp.Status)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}
	return nil
}

func httpPost(uri string, header *url.Values, dataMap map[string]string, response interface{}) error {
	log.Printf("URI %s", uri)
	data, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Errorf("failed to marshal request data to JSON: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if header != nil {
		for k, vv := range *header {
			for _, v := range vv {
				req.Header.Set(k, v)
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("expected HTTP 2xx/3xx, got %s", resp.Status)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}
	return nil
}