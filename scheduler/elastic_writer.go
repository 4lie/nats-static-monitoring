package scheduler

import (
	"fmt"

	"github.com/4lie/nats-static-monitoring/monitor"
	"github.com/olivere/elastic"
)

type ElasticWriter struct {
	ElasticClient *elastic.Client
}

func NewElasticWriter(elasticClient *elastic.Client) *ElasticWriter {
	return &ElasticWriter{
		ElasticClient: elasticClient,
	}
}

func (e ElasticWriter) Write(response map[string]*monitor.Response) {
	for _, value := range response {
		id := fmt.Sprintf("%s-%s", value.Index, value.Key)
		e.ElasticClient.Index().Index(value.Type).Type("doc").Id(id).BodyJson(value.Body)
	}
}
