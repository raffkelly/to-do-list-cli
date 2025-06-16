package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (config *Config) doAuthenticatedRequest(client *http.Client, req *http.Request, result interface{}) error {
	config.attachJWT(req)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = decodeResponse(res, result)
	//handle jwt refreshing here now
	if err != nil && err.Error() == "unable to verify JWT" {
		//attempt refresh
		_, err = config.refreshUserLogin(config.CurrentUser.Username)
		if err != nil {
			return err
		}

		newReq, err := cloneRequest(req)
		if err != nil {
			return err
		}

		config.attachJWT(newReq)
		res, err = client.Do(newReq)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return decodeResponse(res, result)
	} else if err != nil {
		return err
	}
	return nil
}

func cloneRequest(req *http.Request) (*http.Request, error) {
	newReq := req.Clone(req.Context())

	if req.Body != nil {
		var bodyBytes []byte
		if req.GetBody != nil {
			rc, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			bodyBytes, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("cannot clone request: Body is not repeatable")
		}
		newReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	return newReq, nil
}

func (config *Config) attachJWT(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+config.CurrentUser.JWT)
}

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
