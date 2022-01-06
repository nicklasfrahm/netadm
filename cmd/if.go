package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var ifCmd = &cobra.Command{
	Use:   "if",
	Short: "List network interfaces",
	Long: `A utility command that allows you to
list all available network interfaces
that can be used for the operations with
this CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get all network interfaces.
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}

		// Check if any interfaces were found.
		if len(interfaces) == 0 {
			return errors.New("no interfaces found")
		}

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "Interface\tMAC Address\tNetwork Addresses\n")

		// Print all interfaces.
		for _, iface := range interfaces {
			// Skip interface if it does not have
			// a MAC address. An examples of this
			// is the loopback interface.
			mac := iface.HardwareAddr.String()
			if mac == "" {
				continue
			}

			// Skip if the interface is not up.
			if iface.Flags&net.FlagUp == 0 {
				continue
			}

			// Skip if no addresses are assigned
			// to the interface.
			addrs, err := iface.Addrs()
			if err != nil {
				return err
			}
			if len(addrs) == 0 {
				continue
			}
			addresses := make([]string, len(addrs))
			for i, addr := range addrs {
				addresses[i] = addr.String()
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", iface.Name, iface.HardwareAddr.String(), strings.Join(addresses, ","))
		}

		if err := w.Flush(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ifCmd)
}
