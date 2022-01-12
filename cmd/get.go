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
	Use:   "get <mac> [key ...]",
	Short: "Read configuration keys",
	Long: `A command that allows you to read the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the MAC address is valid.
		mac, err := net.ParseMAC(args[0])
		if err != nil {
			return err
		}

		// Check if all keys are valid.
		keys := args[1:]
		for i, key := range keys {
			// Normalize key name.
			keys[i] = strings.ToLower(key)

			// Check if key is valid.
			if nsdp.RecordTypeNames[key] == nil {
				return fmt.Errorf(`unknown configuration key "%s"`, key)
			}
		}

		// Create slice to hold results.
		messages := make([]nsdp.Message, 0)

		// Retry operation if retries is greater than 0.
		for i := uint(0); i <= retries; i++ {
			// Create new message.
			request := nsdp.NewMessage(nsdp.ReadRequest)

			// Add request records.
			for _, key := range keys {
				request.Records = append(request.Records, nsdp.Record{
					Type: nsdp.RecordTypeNames[key].ID,
				})
			}

			// Create context to handle timeout.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Run scan for devices.
			msgs, err := nsdp.RequestMessages(interfaceName, request,
				nsdp.WithContext(ctx),
				nsdp.WithMAC(&mac),
			)
			if err != nil {
				return err
			}

			// Deduplicate results from all attempts.
			messages = nsdp.DeduplicateMessages(messages, msgs)
		}

		// Check if any devices were found.
		if len(messages) == 0 {
			return errors.New("no switches found")
		}

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)

		// Fetch table columns from first message.
		for _, record := range messages[0].Records {
			fmt.Fprintf(w, "%s\t", nsdp.RecordTypeIDs[record.Type].Name)
		}
		fmt.Fprintln(w)

		// Print requested columns.
		for _, message := range messages {
			for _, record := range message.Records {
				// Get corresponding record type.
				recordType := nsdp.RecordTypeIDs[record.Type]
				if recordType == nil {
					fmt.Fprint(w, "\t")
					continue
				}

				// Parse values into their according types.
				var value reflect.Value
				switch recordType.Example.(type) {
				case string:
					value = reflect.ValueOf(string(record.Value))
				case bool:
					value = reflect.ValueOf(record.Value[0] == 1)
				case net.HardwareAddr:
					value = reflect.ValueOf(net.HardwareAddr(record.Value))
				case net.IP:
					value = reflect.ValueOf(net.IP(record.Value))
				}

				fmt.Fprintf(w, "%v\t", value)
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
