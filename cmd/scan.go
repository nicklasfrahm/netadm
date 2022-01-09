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

var retries uint

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for devices",
	Long: `Scan for Netgear switches in your local network.

The program will exit with a non-zero exit
code if the scan does not return anything.

If your scan doesn't return any devices
depite them being present on your network
please increase the timeout and try again.

Another known issue is that the scan does
not return a device if the device needs to
refresh its ARP cache by performing a MAC
address lookup of the host based on its IP.
This happens on the very first interaction
or when the cache naturally expires, which
appears to be every 5 minutes or so. This
issue can be avoided using the retries flag.`,
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

			// If more devices were found, update result.
			if len(devs) > len(devices) {
				devices = devs
			}
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
	scanCmd.Flags().UintVarP(&retries, "retries", "r", 0, "number of retries to perform")
	scanCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(scanCmd)
}
