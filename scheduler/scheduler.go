package scheduler

import (
	"time"

	"github.com/4lie/nats-static-monitoring/config"
	"github.com/4lie/nats-static-monitoring/monitor"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	ElasticWriter *ElasticWriter
}

func New(writer *ElasticWriter) Scheduler {
	return Scheduler{
		ElasticWriter: writer,
	}
}

func (s *Scheduler) fetchAndReport(monitorServers []config.MonitorServer) {
	clients := monitor.Init(monitorServers)
	for _, client := range clients {
		response, err := client.GetStats()
		if err != nil {
			logrus.Errorf("unable to fetch data from %s: %s", client.URL, err)
		}

		s.ElasticWriter.Write(response)
	}
}

func (s *Scheduler) Start(cronPattern string, monitorServers []config.MonitorServer) error {
	c := cron.New()

	err := c.AddFunc(cronPattern, func() {
		logrus.Infof("scheduler start at: %s", time.Now().Format(time.RFC3339))
		s.fetchAndReport(monitorServers)
	})
	if err != nil {
		return err
	}

	c.Start()

	return nil
}
