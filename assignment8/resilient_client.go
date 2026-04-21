package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

// Constants as per requirements
const (
	MaxRetries = 5
	BaseDelay  = 500 * time.Millisecond
)

// IsRetryable returns true if the error is temporary (network timeout, codes 429, 500, 502, 503, 504).
func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		// Check for network timeout
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
		return false
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}

	return false
}

// CalculateBackoff implements exponential backoff with added Full Jitter.
func CalculateBackoff(attempt int) time.Duration {
	// attempt starts from 1
	// Exponential backoff part: base * 2^(attempt-1)
	exponent := 1 << (attempt - 1)
	temp := BaseDelay * time.Duration(exponent)
	
	// Full Jitter: random value between 0 and temp
	// For attempt 1, temp = 500ms, jittered range [0, 500ms)
	// Seed is important for randomization
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Int63n(int64(temp)))
}

// ExecutePayment handles the retry logic.
func ExecutePayment(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Use a new request object with context for each attempt
		// req.Clone(ctx) is better in Go 1.13+
		reqWithCtx := req.Clone(ctx)
		
		resp, err := client.Do(reqWithCtx)
		
		if !IsRetryable(resp, err) {
			// If not retryable, return immediately (success or non-retryable error/4xx)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		// If it's retryable
		if attempt < MaxRetries {
			backoff := CalculateBackoff(attempt)
			fmt.Printf("Attempt %d failed: waiting %v...\n", attempt, backoff)
			
			if resp != nil {
				resp.Body.Close()
			}

			select {
			case <-time.After(backoff):
				// continue to next attempt
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		} else {
			// Last attempt failed
			return resp, err
		}
	}

	return nil, fmt.Errorf("max retries exceeded")
}

func main() {
	fmt.Println("=== Starting Resilient HTTP Client Simulation ===")
	
	// Seed once
	rand.Seed(time.Now().UnixNano())

	// Create a test server
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		fmt.Printf("[Server] Received request #%d\n", requestCount)
		if requestCount <= 3 {
			// Returns 503 Service Unavailable for the first 3 requests
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Returns 200 OK and JSON for the 4th request
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "success"}`)
	}))
	defer server.Close()

	fmt.Printf("[Client] Server URL: %s\n", server.URL)
	
	// Configure client (constants already defined globally)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequest("POST", server.URL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	fmt.Println("[Client] Executing payment...")
	resp, err := ExecutePayment(ctx, req)
	if err != nil {
		fmt.Printf("Final result: Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Attempt %d: Success!\n", requestCount)
		fmt.Printf("Response: %s", string(body))
	} else {
		fmt.Printf("Final status code: %d\n", resp.StatusCode)
	}
}
