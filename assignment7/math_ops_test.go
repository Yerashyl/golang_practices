package assignment7_test

import (
	"testing"

	"assignment7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		want        int
		wantErr     bool
		errContains string
	}{
		{
			name:    "success",
			a:       10,
			b:       2,
			want:    5,
			wantErr: false,
		},
		{
			name:        "division by zero",
			a:       10,
			b:       0,
			want:    0,
			wantErr: true,
			errContains: "division by zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assignment7.Divide(tt.a, tt.b)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{
			name: "both positive",
			a:    10,
			b:    5,
			want: 5,
		},
		{
			name: "positive minus zero",
			a:    10,
			b:    0,
			want: 10,
		},
		{
			name: "negative minus positive",
			a:    -10,
			b:    5,
			want: -15,
		},
		{
			name: "both negative",
			a:    -10,
			b:    -5,
			want: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assignment7.Subtract(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
