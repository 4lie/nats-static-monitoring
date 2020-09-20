package scheduler

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

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
		ctx := context.Background()

		index := fmt.Sprintf("%s-%s", value.Type, value.Index)

		_, err := e.ElasticClient.Index().Index(index).Type("doc").Id(value.Key).BodyJson(value.Body).Do(ctx)
		if err != nil {
			logrus.Errorf("unable to write data to elastic index %s with key %s: %s", index, value.Key, err.Error())
		}

		_, err = e.ElasticClient.Forcemerge(index).Do(ctx)
		if err != nil {
			logrus.Errorf("unable to fore merge data for elastic index %s: %s", index, err.Error())
		}
	}
}
