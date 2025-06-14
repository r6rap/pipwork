package checker

import (
	"net/http"
	"time"
)

type HTTPResult struct {
	Status  string         // UP or DOWN
	Latency time.Duration
	StatusCode int
	Error   string
}

var defaultTransport = &http.Transport{
	MaxIdleConns: 100,
	IdleConnTimeout: 90 * time.Second,
	DisableCompression: false,
}

var httpClient = &http.Client{
	Timeout: 7 * time.Second,
	Transport: defaultTransport,
}

func CheckHTTP(url string) HTTPResult {
	// save time when checking
	start := time.Now()

	// send a get request
	resp, err := httpClient.Get(url)
	// response time received
	latency := time.Since(start)

	if err != nil {
		return HTTPResult{Status: "DOWN", Latency: latency, Error: err.Error()}
	}
	// close response body after func finished
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return HTTPResult{Status: "UP", Latency: latency, StatusCode: http.StatusOK}
	}

	return HTTPResult{Status: "DOWN", Latency: latency, Error: resp.Status}
}