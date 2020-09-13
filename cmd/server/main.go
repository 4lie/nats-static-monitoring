package server

import (
	"github.com/4lie/nats-static-monitoring/config"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
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
