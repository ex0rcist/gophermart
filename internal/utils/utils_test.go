package utils

import (
	"net/http"
	"regexp"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func TestIntToDuration(t *testing.T) {
	tests := []struct {
		input    int
		expected time.Duration
	}{
		{input: 0, expected: 0 * time.Second},
		{input: 1, expected: 1 * time.Second},
		{input: 60, expected: 60 * time.Second},
		{input: -1, expected: -1 * time.Second},
		{input: 3600, expected: 3600 * time.Second},
	}

	for _, tt := range tests {
		result := IntToDuration(tt.input)
		if result != tt.expected {
			t.Errorf("IntToDuration(%d) = %v; expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestHeadersToStr(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected string
	}{
		{
			name: "Single header with single value",
			headers: http.Header{
				"Content-Type": {"application/json"},
			},
			expected: "Content-Type:application/json",
		},
		{
			name: "Single header with multiple values",
			headers: http.Header{
				"Accept": {"text/plain", "text/html"},
			},
			expected: "Accept:text/html, Accept:text/plain",
		},
		{
			name: "Multiple headers with single values",
			headers: http.Header{
				"Content-Type": {"application/json"},
				"User-Agent":   {"Go-http-client/1.1"},
			},
			expected: "Content-Type:application/json, User-Agent:Go-http-client/1.1",
		},
		{
			name:     "Empty headers",
			headers:  http.Header{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HeadersToStr(tt.headers)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	regex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

	requestID := GenerateRequestID()
	if !regex.MatchString(requestID) {
		t.Errorf("GenerateRequestID() returned invalid UUIDv4: %s", requestID)
	}
}

func TestLuhnCheck(t *testing.T) {
	type args struct {
		nums string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"0 is invalid", args{"0"}, false},
		{"1 is invalid", args{"1"}, false},
		{"11 is invalid", args{"11"}, false},
		{"70483 is invalid", args{"70483"}, false},
		{"349926205465199 is invalid", args{"349926205465199"}, false},
		{"00 is valid", args{"00"}, true},
		{"18 is valid", args{"18"}, true},
		{"70482 is valid", args{"70482"}, true},
		{"349926205465194 is valid", args{"349926205465194"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnCheck(tt.args.nums); got != tt.want {
				t.Errorf("checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 10
	randomString := GenerateRandomString(length)

	if utf8.RuneCountInString(randomString) == 0 {
		t.Errorf("generated string is empty")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("hashed password is empty")
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	err = ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("password comparison failed: %v", err)
	}

	err = ComparePassword(hashedPassword, "wrongpassword")
	if err == nil {
		t.Errorf("expected comparison to fail, but it succeeded")
	}
}

type LuhnTestStruct struct {
	CardNumber string `validate:"luhn"`
}

func TestLuhnValidation(t *testing.T) {
	validate := validator.New()

	err := validate.RegisterValidation("luhn", luhnValidation)
	if err != nil {
		t.Fatalf("Failed to register luhn validation: %v", err)
	}

	validNumber := LuhnTestStruct{CardNumber: "79927398713"}
	err = validate.Struct(validNumber)
	if err != nil {
		t.Errorf("expected valid Luhn number, got error: %v", err)
	}

	invalidNumber := LuhnTestStruct{CardNumber: "12345678901"}
	err = validate.Struct(invalidNumber)
	if err == nil {
		t.Errorf("expected invalid Luhn number, but got no error")
	}
}
