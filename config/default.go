package config

const Namespace = "nsm"

//nolint:lll
const Default = `
logger:
  level: debug

monitor-servers:
- server: http://127.0.0.1:8222
  type: NATS
  connect-timeout: 1s 
  endpoint-uris:
  - /varz
  - /connz?subs=detail&auth=1&state=any
  - /routez?sub=detail
  - /gatewayz?accs=1
  - /leafz?subs=1
  - /subsz?subs=1
- server: http://127.0.0.1:8223
  type: NATS-Streaming
  connect-timeout: 1s 
  endpoint-uris:
  - /varz
  - /connz?subs=detail&auth=1&state=any
  - /routez?sub=detail
  - /gatewayz?accs=1
  - /leafz?subs=1
  - /streaming/serverz
  - /streaming/storez
  - /streaming/clientsz?subs=1
  - /streaming/channelsz?subs=1

elasticsearch:
  servers:
    - 127.0.0.1:9200
`
