package cmd

import (
	"os"

	"github.com/nicklasfrahm/netadm/pkg/fmt"
	"github.com/nicklasfrahm/netadm/pkg/nsdp"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for devices",
	Long: `Scan for devices in your local network.

The program will exit with a non-zero exit
code if the scan does not return anything.

If your scan doesn't return any devices
despite them being present on your network
please increase the timeout and try again.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := "all"
		keys := []string{"name", "model", "mac", "ip", "dhcp", "firmware", "encryptionmode"}

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
	scanCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	scanCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(scanCmd)
}
