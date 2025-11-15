package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/complytime/gemara-mcp-server/mcp"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/complytime/gemara-mcp-server/version"
)

var (
	configFile  string
	transport   string
	logFilePath string
)

var rootCmd = &cobra.Command{
	Use:   "gemara-mcp-server",
	Short: "Gemara CUE MCP Server",
	Long:  "A Model Context Protocol server for Gemara (GRC Engineering Model for Automated Risk Assessment) with CUE support",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mcp.Config{
			Version:     version.GetVersion(),
			LogFilePath: logFilePath,
		}

		server, err := mcp.NewServer(cfg)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		return server.Start()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gemara CUE MCP Server %s\n", version.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().StringVar(&configFile, "config", "", "path to configuration file (not currently used)")
	rootCmd.Flags().StringVar(&transport, "transport", "stdio", "transport mode (stdio/sse) - only stdio is currently supported")
	rootCmd.Flags().StringVar(&logFilePath, "log-file", "", "path to log file (default: stderr)")

	// Bridge glog flags with pflag for cobra compatibility
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func main() {
	defer glog.Flush()

	if err := rootCmd.Execute(); err != nil {
		glog.Errorf("Command execution failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
