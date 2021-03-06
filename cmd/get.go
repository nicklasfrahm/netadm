package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
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
		selector := nsdp.NewSelector()
		// Allow usage of keyword "all" to select all devices.
		if args[0] != "all" {
			// Check if the device is identified via its IP address.
			ip := net.ParseIP(args[0])
			if ip == nil {
				// Fall back to MAC address device identification.
				mac, err := net.ParseMAC(args[0])
				if err != nil {
					return errors.New("device identifier must be a MAC address or an IP address")
				}
				selector.SetMAC(&mac)
			} else {
				selector.SetIP(&ip)
			}
		}
		// Check if all keys are valid.
		keys := args[1:]
		for i, key := range keys {
			// Normalize key name.
			keys[i] = strings.ToLower(key)

			// Check if key is valid.
			if nsdp.RecordTypeByName[key] == nil {
				return fmt.Errorf(`unknown configuration key "%s"`, key)
			}
		}

		// Create slice to hold results.
		devices := make([]nsdp.Device, 0)

		// Retry operation if retries is greater than 0.
		for i := uint(0); i <= retries; i++ {
			// Create new message.
			request := nsdp.NewMessage(nsdp.ReadRequest)

			// Add request records.
			for _, key := range keys {
				request.Records = append(request.Records, nsdp.Record{
					ID: nsdp.RecordTypeByName[key].ID,
				})
			}

			// Create context to handle timeout.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Run scan for devices.
			devs, err := nsdp.RequestDevices(interfaceName, request,
				nsdp.WithContext(ctx),
				nsdp.WithSelector(selector),
			)
			if err != nil {
				return err
			}

			// Deduplicate results from all attempts.
			devices = nsdp.DeduplicateDevices(devices, devs)
		}

		// Check if any devices were found.
		if len(devices) == 0 {
			return errors.New("no switches found")
		}

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)

		// Fetch table columns from desired keys.
		for _, key := range keys {
			// Print column header.
			rt := nsdp.RecordTypeByName[key]
			fmt.Fprintf(w, "%s\t", strings.ToUpper(rt.Name))
		}
		fmt.Fprintln(w)

		// Print requested keys for each device. Note that we
		// unmarshal the message into a Device because this
		// allows it to easily group messages of the same type.
		for _, device := range devices {
			// Print the desired columns.
			for _, key := range keys {
				// Fetch field from device.
				name := nsdp.RecordTypeByName[key].Name
				field := reflect.ValueOf(device).FieldByName(name)
				if field.IsValid() {
					// Print field value.
					fmt.Fprintf(w, "%v\t", field)
				} else {
					// This happens if the field is a known message
					// type but not defined inside the Device struct.
					// Make sure to add the according field to the
					// Device struct to prevent this.
					fmt.Fprintf(w, "n/a\t")
				}
			}

			fmt.Fprintln(w)
		}

		return w.Flush()
	},
}

func init() {
	getCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	getCmd.MarkFlagRequired("interface")

	rootCmd.AddCommand(getCmd)
}
