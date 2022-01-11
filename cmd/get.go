package cmd

import (
	"errors"
	"fmt"
	"net"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <mac> [key ...]",
	Short: "Read configuration keys",
	Long: `A command that allows you to read the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the MAC address is valid.
		_, err := net.ParseMAC(args[0])
		if err != nil {
			return err
		}

		// Check if all keys are valid.
		keys := args[1:]
		for _, key := range keys {
			if nsdp.RecordNames[key] == nil {
				return fmt.Errorf(`unknown configuration key "%s"`, key)
			}
		}

		// TODO: Implement "get" command logic.
		return errors.New("not implemented")

		// devices := make([]nsdp.Device, 0)

		// // Create new message.
		// request := nsdp.NewMessage(nsdp.ReadRequest)

		// // TODO: Add records according to requested keys.

		// // Retry operation if retries is greater than 0.
		// for i := uint(0); i <= retries; i++ {
		// 	// Create context to handle timeout.
		// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
		// 	defer cancel()

		// 	// Run scan for devices.
		// 	devs, err := nsdp.RequestDevice(interfaceName, request,
		// 		nsdp.WithContext(ctx),
		// 		nsdp.WithMAC(&mac),
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	// Deduplicate results from all attempts.
		// 	devices = nsdp.Deduplicate(devices, devs)
		// }

		// // Check if any devices were found.
		// if len(devices) == 0 {
		// 	return errors.New("no switches found")
		// }

		// // Create table with tabwriter.
		// w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		// fmt.Fprintf(w, "Name\tModel\tMAC Address\tIP Address\tDHCP\tFirmware\n")

		// // Print simple list of switches.
		// for _, device := range devices {
		// 	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%s\n", device.Name, device.Model, device.MAC.String(), device.IP.String(), device.DHCP, device.Firmware)
		// }

		// return w.Flush()
	},
}

func init() {
	getCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	getCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(getCmd)
}
