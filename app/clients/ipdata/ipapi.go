package ipdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type client struct {
	requestTimeout time.Duration
}

func NewClient(requestTimeout time.Duration) *client {
	return &client{
		requestTimeout: requestTimeout,
	}
}

func (c *client) GetRequestLocation(ip string) (string, error) {
	r, err := http.NewRequest(`GET`, fmt.Sprintf(`https://ipapi.co/%s/json/`, ip), nil)
	if err != nil {
		return ``, err
	}
	client := &http.Client{
		Timeout: c.requestTimeout,
	}
	resp, err := client.Do(r)
	if err != nil {
		return ``, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return ``, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ``, fmt.Errorf("got non OK response: %s", string(response))
	}

	responseData := struct {
		CountryCode string `json:"country_code"`
		Error       bool   `json:"error"`
		Reason      string `json:"reason"`
	}{}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return ``, err
	}
	if responseData.Error {
		return ``, errors.New(responseData.Reason)
	}
	return responseData.CountryCode, nil
}
