package accrual

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetBonuses_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"order":"12345","status":"PROCESSED","accrual":150.50}`))
		assert.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.Equal(t, "12345", response.OrderNumber)
	assert.Equal(t, StatusProcessed, response.Status)
	accrualAmount, _ := decimal.NewFromString("150.50")
	assert.Equal(t, accrualAmount, response.Amount)
}

func TestGetBonuses_NoContent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "12345", response.OrderNumber)
	assert.Equal(t, StatusRegistered, response.Status)
	assert.Equal(t, decimal.NewFromInt(0), response.Amount)
}

func TestGetBonuses_TooManyRequests(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	fmt.Println(err)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusTooManyRequests, err.HTTPStatus)
	assert.Equal(t, 5*time.Second, err.RetryAfter)
}

func TestGetBonuses_UnexpectedStatus(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
}

func TestGetBonuses_RequestCreationError(t *testing.T) {
	client := NewClient("http://%41:8080", 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	assert.Nil(t, response)
	assert.NotNil(t, err)
}

func TestGetBonuses_BodyReadError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"order":`))
		assert.Nil(t, err)
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	ctx := context.Background()

	response, err := client.GetBonuses(ctx, "12345")

	assert.Nil(t, response)
	assert.NotNil(t, err)
}
