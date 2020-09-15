package cmd

import (
	"os"

	"github.com/4lie/nats-static-monitoring/cmd/scheduler"

	"github.com/4lie/nats-static-monitoring/cmd/server"
	"github.com/4lie/nats-static-monitoring/config"
	"github.com/4lie/nats-static-monitoring/log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.Init("config.yaml")

	log.SetupLogger(cfg.Logger)

	root := &cobra.Command{
		Use:   "nsm",
		Short: "NATS / NATS Streaming Static Monitoring Service",
	}

	server.Register(root, cfg)
	scheduler.Register(root, cfg)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
