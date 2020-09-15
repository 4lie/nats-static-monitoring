package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/4lie/nats-static-monitoring/config"

	"github.com/go-resty/resty/v2"
)

type Response struct {
	Body  json.RawMessage
	Key   string
	Index string
	Type  string
}

type Client struct {
	URL       string
	Type      string
	Alias     string
	Endpoints []string
	client    *resty.Client
}

func Init(cfgs []config.MonitorServer) []*Client {
	clients := make([]*Client, len(cfgs))
	for index, cfg := range cfgs {
		clients[index] = New(cfg.Server, cfg.ConnectTimeout, cfg.Type, cfg.Alias, cfg.Endpoints)
	}

	return clients
}

func New(url string, timeout time.Duration, clientType, alias string, uris []string) *Client {
	return &Client{
		URL:       url,
		client:    resty.New().SetHostURL(url).SetTimeout(timeout),
		Endpoints: uris,
		Alias:     alias,
		Type:      clientType,
	}
}

//nolint:goerr113
func (c *Client) GetStats() (map[string]*Response, error) {
	response := make(map[string]*Response)

	for _, URI := range c.Endpoints {
		resp, err := c.client.R().Get(URI)
		if err != nil {
			return nil, fmt.Errorf("error sending request to %s: %w", URI, err)
		}

		if resp.StatusCode() != http.StatusOK {
			return nil,
				fmt.Errorf("request for %s failed with status = %d and error = %s", URI, resp.StatusCode(), resp.String())
		}

		rawData := make(map[string]interface{})

		err = json.Unmarshal(resp.Body(), &rawData)
		if err != nil {
			return nil, fmt.Errorf("error parsing request body: %w", err)
		}

		index := c.extractIndex(URI)
		monitorResponse := new(Response)
		monitorResponse.Type = c.Type
		monitorResponse.Body = resp.Body()
		monitorResponse.Key = c.Alias
		monitorResponse.Index = index
		response[index] = monitorResponse
	}

	return response, nil
}

func (c *Client) extractIndex(uri string) string {
	splitURI := strings.SplitAfter(uri, "/")
	uri = strings.TrimLeft(splitURI[len(splitURI)-1], "/")
	endpoint := strings.SplitN(uri, "?", 2)

	return endpoint[0]
}
