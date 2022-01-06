package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var interfaceName string
var timeout time.Duration
var help bool

var rootCmd = &cobra.Command{
	Use:   "nsdp",
	Short: "CLI for the Netgear Switch Discovery Protocol (NSDP)",
	Long: `A command line interface to manage Netgear Smart Switches
via the UDP-based Netgear Switch Discovery Protocol (NSDP).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", time.Second, "timeout for commands")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
