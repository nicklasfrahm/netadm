package cmd

import (
	"os"

	"github.com/nicklasfrahm/netadm/pkg/fmt"
	"github.com/nicklasfrahm/netadm/pkg/nsdp"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <device> [key ...]",
	Short: "Read configuration keys",
	Long: `A command that allows you to read the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		keys := args[1:]

		devices, err := nsdp.Get(id, keys,
			nsdp.WithInterfaceName(interfaceName),
			nsdp.WithRetries(retries),
			nsdp.WithTimeout(timeout),
		)
		if err != nil {
			return err
		}

		// Print results.
		fmt.Table(os.Stdout, devices, keys)

		return nil
	},
}

func init() {
	getCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	getCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(getCmd)
}
