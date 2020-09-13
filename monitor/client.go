package monitor

import "encoding/json"

type Client interface {
	GetStats() (json.RawMessage, error)
}
