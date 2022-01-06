package cmd

import (
	"errors"
	"fmt"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for devices",
	Long: `Scan for Netgear switches in your local network.

The program will exit with a non-zero exit
code if the scan does not return anything.

If your scan doesn't return any devices
depite them being present on your network
please increase the timeout and try again.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		devices, err := nsdp.Scan(interfaceName,
			nsdp.Timeout(timeout),
		)
		if err != nil {
			return err
		}

		// Check if any devices were found.
		if len(*devices) == 0 {
			return errors.New("no switches found")
		}

		// Print simple list of switches.
		for _, device := range *devices {
			fmt.Println(device.Name)
		}

		return nil
	},
}

func init() {
	scanCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	scanCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(scanCmd)
}
