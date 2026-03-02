package analyzer

import (
	"testing"
)

//1: startsWithLower

func TestStartsWithLower(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"server started", true},
		{"Server started", false},
		{" server started", true},
		{"   Server started", false},
		{"", true},    // пустое считаем ок
		{"   ", true}, // только пробелы ок
		{"123 started", true},
	}

	for _, tt := range tests {
		got := startsWithLower(tt.in)
		if got != tt.want {
			t.Fatalf("startsWithLower(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

//2: isOnlyEngLetters

func TestIsOnlyEngLetters(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"server started", true},
		{"Server started", true},
		{"сервер запущен", false},
		{"server запущен", false},
		{"hello мир", false},
	}

	for _, tt := range tests {
		got := isOnlyEngLetters(tt.in)
		if got != tt.want {
			t.Fatalf("isOnlyEngLetters(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

//3: isNotSpecSymbols

func TestIsNotSpecSymbols(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"server started", true},
		{"server started 123", true},
		{"server-started", false},
		{"server_started", false},
		{"server started!", false},
		{"warning: server started", false},
		{"server 🚀", false},
	}

	for _, tt := range tests {
		got := isNotSpecSymbols(tt.in)
		if got != tt.want {
			t.Fatalf("isNotSpecSymbols(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

//4: isNotPrivateInfo

func TestIsNotPrivateInfo(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"user authenticated", true},
		{"token expired", true},
		{"invalid token", true},
		{"token: 123", false},
		{"token=123", false},
		{"password: 123", false},
		{"password is 123", false},
		{"api_key=abc", false},
		{"apikey: abc", false},
		{"token: %s", false},
		{"password %v", false},
	}

	for _, tt := range tests {
		got := isNotPrivateInfo(tt.in)
		if got != tt.want {
			t.Fatalf("isNotPrivateInfo(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
