package monitor

import (
	"encoding/json"

	"github.com/4lie/nats-static-monitoring/config"
)

type Response struct {
	Body     json.RawMessage
	ServerID string
	Type     string
}

type Client interface {
	GetStats() (map[string]*Response, error)
	Type() string
}

func Init(cfgs []config.MonitorServer) []Client {
	clients := make([]Client, len(cfgs))
	for index, cfg := range cfgs {
		clients[index] = NewNATSClient(cfg.Server, cfg.ConnectTimeout, cfg.Type, cfg.EndpointURIs)
	}

	return clients
}
