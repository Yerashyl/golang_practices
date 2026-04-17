package assignment7_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"assignment7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRate(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		from, to       string
		wantRate       float64
		wantErr        bool
		errContains    string
	}{
		{
			name: "Successfull scenario",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"base":"USD", "target":"EUR", "rate":0.85}`)
			},
			from:     "USD",
			to:       "EUR",
			wantRate: 0.85,
			wantErr:  false,
		},
		{
			name: "API Business Error (400)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"error":"invalid currency pair"}`)
			},
			wantErr:     true,
			errContains: "api error: invalid currency pair",
		},
		{
			name: "API Business Error (404)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"error":"invalid currency pair"}`)
			},
			wantErr:     true,
			errContains: "api error: invalid currency pair",
		},
		{
			name: "Malformed JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"base": "USD", "rate": "invalid"}`)
			},
			wantErr:     true,
			errContains: "decode error",
		},
		{
			name: "Slow Response/Timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(6 * time.Second)
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"rate":0.85}`)
			},
			wantErr:     true,
			errContains: "network error",
		},
		{
			name: "Server Panic / 500 Internal Server Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"error":"panic happened"}`)
			},
			wantErr:     true,
			errContains: "api error: panic happened",
		},
		{
			name: "Empty Body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr:     true,
			errContains: "decode error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			service := assignment7.NewExchangeService(server.URL)
			gotRate, err := service.GetRate(tt.from, tt.to)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantRate, gotRate)
			}
		})
	}
}
