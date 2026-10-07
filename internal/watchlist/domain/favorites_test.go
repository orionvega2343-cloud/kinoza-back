package domain

import (
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "watched", status: "watched", want: true},
		{name: "watching", status: "watching", want: true},
		{name: "planned", status: "planned", want: true},
		{name: "invalid", status: "invalid", want: false},
		{name: "empty", status: "", want: false},
		{name: "uppercase", status: "Planned", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
