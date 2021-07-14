package app

import (
	"fmt"

	"istio.io/pkg/log"

	"github.com/spf13/cobra"

	"github.com/zirain/mcpoverxds/pkg/bootstrap"
	"github.com/zirain/mcpoverxds/pkg/cmd"
)

var (
	loggingOptions *log.Options
	serverArgs     *bootstrap.ServerArgs
)

func NewRootCommand() *cobra.Command {
	loggingOptions = log.DefaultOptions()

	rootCmd := &cobra.Command{
		Use:          "mcp-over-xds",
		Short:        "mcp-over-xds",
		Long:         "This is a demo for istio mcp-over-xds usage.",
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			cmd.AddFlags(c)
			return nil
		},
	}
	loggingOptions.AttachCobraFlags(rootCmd)

	serverCmd := newServerCommand()
	addServerFlags(serverCmd)
	rootCmd.AddCommand(serverCmd)

	return rootCmd
}

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start mcp-over-xds server.",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(c *cobra.Command, args []string) error {
			if err := log.Configure(loggingOptions); err != nil {
				return err
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			cmd.PrintFlags(c.Flags())

			// Create the stop channel for all of the servers.
			stop := make(chan struct{})

			// Create the server for the discovery service.
			xdsServer, err := bootstrap.NewServer(serverArgs)
			if err != nil {
				return fmt.Errorf("failed to create discovery service: %v", err)
			}

			// Start the server
			if err := xdsServer.Start(stop); err != nil {
				return fmt.Errorf("failed to start discovery service: %v", err)
			}

			cmd.WaitSignal(stop)
			return nil
		},
	}
}

func addServerFlags(c *cobra.Command) {
	serverArgs = &bootstrap.ServerArgs{}

	// Process commandline args.
	c.PersistentFlags().StringVarP(&serverArgs.DiscoveryAddress, "discovery", "d",
		"istiod.istio-system.svc:15010",
		"The address of istiod xds server, default values is istiod.istio-system.svc:15010")
}
