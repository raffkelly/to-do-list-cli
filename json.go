package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func decodeResponse(resp *http.Response, successPtr interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return json.Unmarshal(body, successPtr)
	}

	// Handle error case: try to unmarshal error payload first
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
		return errors.New(errResp.Error)
	}
	// fallback: unknown error format
	return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
}
