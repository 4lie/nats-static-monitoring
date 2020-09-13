package config

const Namespace = "nsm"

//nolint:lll
const Default = `
logger:
  level: debug

monitor-servers:
- server: http://demo.nats.io:8222
  type: NATS
  connect-timeout: 1s 
  endpoint-uris:
  - /varz
  - /connz?subs=detail

elasticsearch:
  servers:
    - 127.0.0.1:9200
`
