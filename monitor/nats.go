package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type NATSClient struct {
	baseURL    string
	clientType string
	URIs       []string
	client     *resty.Client
}

func NewNATSClient(url string, timeout time.Duration, clientType string, uris []string) *NATSClient {
	return &NATSClient{
		baseURL:    url,
		client:     resty.New().SetHostURL(url).SetTimeout(timeout),
		URIs:       uris,
		clientType: clientType,
	}
}

//nolint:goerr113
func (c *NATSClient) GetStats() (map[string]*Response, error) {
	response := make(map[string]*Response)

	for _, URI := range c.URIs {
		resp, err := c.client.R().Get(URI)
		if err != nil {
			return nil, fmt.Errorf("error sending request to %s: %w", URI, err)
		}

		if resp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("request failed with status = %d and error = %s", resp.StatusCode(), resp.String())
		}

		rawData := make(map[string]interface{})

		err = json.Unmarshal(resp.Body(), &rawData)
		if err != nil {
			return nil, fmt.Errorf("error parsing request body: %w", err)
		}

		monitorResponse := new(Response)
		monitorResponse.Type = c.clientType
		monitorResponse.Body = resp.Body()
		monitorResponse.ServerID = rawData["server_id"].(string)

		response[c.extractEndpoint(URI)] = monitorResponse
	}

	return response, nil
}

func (c *NATSClient) Type() string {
	return c.clientType
}

func (c *NATSClient) extractEndpoint(uri string) string {
	uri = strings.TrimLeft(uri, "/")
	endpoint := strings.SplitN(uri, "?", 2)

	return endpoint[0]
}
