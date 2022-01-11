package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

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
		devices := make([]nsdp.Device, 0)

		// Retry operation if retries is greater than 0.
		for i := uint(0); i <= retries; i++ {
			// Create context to handle timeout.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Run scan for devices.
			devs, err := nsdp.Scan(interfaceName, nsdp.WithContext(ctx))
			if err != nil {
				return err
			}

			// Deduplicate results from all attempts.
			devices = nsdp.Deduplicate(devices, devs)
		}

		// Check if any devices were found.
		if len(devices) == 0 {
			return errors.New("no switches found")
		}

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "Name\tModel\tMAC Address\tIP Address\tDHCP\tFirmware\n")

		// Print simple list of switches.
		for _, device := range devices {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%s\n", device.Name, device.Model, device.MAC.String(), device.IP.String(), device.DHCP, device.Firmware)
		}

		return w.Flush()
	},
}

func init() {
	scanCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	scanCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(scanCmd)
}
