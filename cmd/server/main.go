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
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Elasticsearch.Servers...),
	}

	elasticClient, err := elastic.NewSimpleClient(opts...)
	if err != nil {
		logrus.Fatalf("failed to create elasticsearch elasticClient: %s", err.Error())
	}

	elasticWriter := scheduler.NewElasticWriter(elasticClient)
	monitoringScheduler := scheduler.New(elasticWriter)

	if err := monitoringScheduler.Start(cfg.Scheduler.CronPattern, cfg.MonitorServers); err != nil {
		logrus.Fatalf("Failed to start scheduler: %s", err.Error())
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logrus.Info("server started!")

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
