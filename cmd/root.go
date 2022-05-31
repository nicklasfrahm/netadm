package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var driver string
var interfaceName string

var timeout time.Duration
var retries uint
var help bool

var rootCmd = &cobra.Command{
	Use:   "netadm",
	Short: "CLI to manage network devices remotely",
	Long: `A command line interface to manage a variety of different
network devices remotely via different protocols.

Note:
  To achieve a consistent behavior all operations
  are executed twice and the results are merged.
  This is done to work around operations that do
  not succeed if the device needs to refresh its
  ARP cache by performing a MAC address lookup of
  the host via the host IP. This happens on the
  the first interaction or, I assume, when the
  cache expires naturally, which appears to be
  every 5 minutes or so.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 100*time.Millisecond, "timeout per attempt")
	rootCmd.PersistentFlags().UintVarP(&retries, "retries", "r", 1, "number of retries to perform")
}

// Execute starts the invocation of the command line interface.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
