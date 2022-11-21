package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jkandasa/file-store/pkg/utils"
)

const DefaultTimeout = time.Second * 30

// ResponseConfig of a request
type ResponseConfig struct {
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"-"`
}

// Client struct
type HttpClient struct {
	httpClient *http.Client
}

// return new client
func newHttpClient(insecure bool, timeout string) *HttpClient {
	var httpClient *http.Client
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient = &http.Client{Transport: customTransport}
	} else {
		httpClient = http.DefaultClient
	}

	timeoutDuration := utils.ToDuration(timeout, DefaultTimeout)
	if timeoutDuration > 0 {
		httpClient.Timeout = timeoutDuration
	} else {
		httpClient.Timeout = DefaultTimeout
	}

	return &HttpClient{httpClient: httpClient}
}

// executeJson execute http request and returns response
func (c *HttpClient) executeJson(url, method string, headers map[string]string, queryParams map[string]interface{},
	body interface{}, responseCode int) (*ResponseConfig, error) {
	// add body, if available
	var buf io.ReadWriter
	if method != http.MethodGet && body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// set as json content
	req.Header.Set("Accept", "application/json")
	// update headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if responseCode > 0 && resp.StatusCode != responseCode {
		return nil, fmt.Errorf("failed with status code. [status: %v, statusCode: %v, body: %s]", resp.Status, resp.StatusCode, string(respBodyBytes))
	}

	respCfg := &ResponseConfig{
		StatusCode: resp.StatusCode,
		URL:        url,
		Method:     method,
		Body:       respBodyBytes,
		Headers:    make(map[string]string),
	}

	// update headers
	for k := range resp.Header {
		respCfg.Headers[k] = resp.Header.Get(k)
	}

	return respCfg, nil
}
