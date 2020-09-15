package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	Timeout      = 2 * time.Second
	Endpoint     = "test-end-point"
	Alias        = "localhost"
	ClientType   = "NATS"
	JSONResponse = `{"server_id":"test_server_id","stats":"everything is fine"}`
)

//nolint:gochecknoglobals
var EndpointURI = []string{"/test/test-end-point?test_param=true"}

type ClientTestSuite struct {
	suite.Suite
	client *Client
}

func newHTTPServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.String() == EndpointURI[0] {
			_, err := rw.Write([]byte(JSONResponse))
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
	}))

	return server
}

func (suite *ClientTestSuite) SetupSuite() {
	server := newHTTPServer()
	suite.client = New(server.URL, Timeout, ClientType, Alias, EndpointURI)
}

func (suite *ClientTestSuite) TestGetStats() {
	r, err := suite.client.GetStats()
	suite.NoError(err)

	natsResponse, ok := r[Endpoint]
	suite.True(ok)

	suite.Equal(Alias, natsResponse.Key)
	suite.Equal(Endpoint, natsResponse.Index)
	suite.Equal(ClientType, natsResponse.Type)
	suite.Equal(JSONResponse, string(natsResponse.Body))
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
