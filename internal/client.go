package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ApiClient struct {
	BaseUrl string
}

func Create(baseUrl string) *ApiClient {
	return &ApiClient{
		BaseUrl: baseUrl,
	}
}

func (c *ApiClient) MakeRequest(method, endpoint string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.BaseUrl + endpoint

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}
