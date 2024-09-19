package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/shopspring/decimal"
)

type AccrualStatus string

const (
	StatusNew        AccrualStatus = "REGISTERED"
	StatusInvalid    AccrualStatus = "INVALID"
	StatusProcessing AccrualStatus = "PROCESSING"
	StatusProcessed  AccrualStatus = "PROCESSED"
)

// клиент для работы с системой начисления бонусов
type Client struct {
	address string
	client  *http.Client
}

type Response struct {
	OrderNumber string          `json:"order"`
	Status      AccrualStatus   `json:"status"`
	Amount      decimal.Decimal `json:"accrual"`
}

type ClientError struct {
	error
	HTTPStatus int
	RetryAfter time.Duration
}

func NewClient(address string, timeout time.Duration) *Client {
	return &Client{
		address: address,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetBonuses(ctx context.Context, orderNumber string) (*Response, *ClientError) {
	url := fmt.Sprintf("http://%s/api/orders/%s", c.address, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ClientError{error: err}
	}

	req.Header.Set("Content-Length", "0")

	logRequest(ctx, url)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, &ClientError{error: err}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &ClientError{error: err}
	}

	logResponse(ctx, res, body)

	// 429 - возвращаем специальную ошибку
	if res.StatusCode == http.StatusTooManyRequests {
		return nil, c.handleErrTooManyRequests(res)
	}

	// 204 - заказ еще не известен бонусной системе;
	// согласно ТЗ, все заказы рано или поздно появятся в accrual,
	// поэтому трактуем как NEW
	if res.StatusCode == http.StatusNoContent {
		logging.LogInfoCtx(ctx, "")
		return &Response{OrderNumber: orderNumber, Status: StatusNew, Amount: decimal.NewFromInt(0)}, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, &ClientError{error: errors.New(http.StatusText(res.StatusCode)), HTTPStatus: res.StatusCode}
	}

	accrualRes := &Response{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(accrualRes)
	if err != nil {
		return nil, &ClientError{error: err}
	}
	return accrualRes, nil
}

func (c *Client) handleErrTooManyRequests(res *http.Response) *ClientError {
	accrualErr := ClientError{
		error:      errors.New(http.StatusText(res.StatusCode)),
		HTTPStatus: res.StatusCode,
	}

	retryAfter, err := strconv.Atoi(res.Header.Get("Retry-After"))
	if err != nil {
		return &ClientError{error: err}
	}
	accrualErr.RetryAfter = utils.IntToDuration(retryAfter)

	return &accrualErr
}

func logRequest(ctx context.Context, url string) {
	logging.LogInfoCtx(ctx, "sending request to: "+url)
}

func logResponse(ctx context.Context, resp *http.Response, respBody []byte) {
	logging.LogDebugCtx(ctx, fmt.Sprintf("response: %v; headers=%s; body=%s", resp.Status, utils.HeadersToStr(resp.Header), respBody))
}
