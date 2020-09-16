package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/4lie/nats-static-monitoring/config"
	"github.com/4lie/nats-static-monitoring/scheduler"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	client, err := elastic.NewClient(elastic.SetURL(cfg.Elasticsearch.Servers...))
	if err != nil {
		logrus.Panicf("failed to create elasticsearch client: %s", err.Error())
	}

	elasticWriter := scheduler.NewElasticWriter(client)
	monitoringScheduler := scheduler.New(elasticWriter)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := monitoringScheduler.Start(cfg.Scheduler.CronPattern, cfg.MonitorServers); err != nil {
			logrus.Fatalf("Failed to start scheduler: %s", err.Error())
		}
	}()

	s := <-sig
	logrus.Infof("signal %s received", s)
}

// Register Server command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "NSM Server Component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
