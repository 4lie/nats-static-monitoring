package config

const Namespace = "nsm"

//nolint:lll
const Default = `
logger:
  level: debug

nats:
  servers:
    - 127.0.0.1:8222

nats-streaming:
  servers:
    - 127.0.0.1:8223

elasticsearch:
  servers:
    - 127.0.0.1:9200
`
