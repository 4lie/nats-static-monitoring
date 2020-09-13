package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type NATSRequest struct {
	Endpoint   string
	Parameters map[string]string
}

type NATSClient struct {
	baseURL string
	client  *resty.Client
}

func NewNATSClient(url string, timeout time.Duration) NATSClient {
	return NATSClient{
		baseURL: url,
		client:  resty.New().SetHostURL(url).SetTimeout(timeout),
	}
}

func (c *NATSClient) GetStats(request NATSRequest) (json.RawMessage, error) {
	resp, err := c.client.R().SetQueryParams(request.Parameters).Get(request.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error sending NATS request to %s: %w", request.Endpoint, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("NATS request failed with status = %d and body = %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
