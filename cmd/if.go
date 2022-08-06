package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"text/tabwriter"

	"github.com/nicklasfrahm/netadm/pkg/nsdp"
	"github.com/spf13/cobra"
)

var ifCmd = &cobra.Command{
	Use:   "if",
	Short: "List network interfaces",
	Long: `A utility command that allows you to
list all available network interfaces
that can be used for the operations with
this CLI.

By default the command will list only
interfaces which are up, have a MAC
and an IPv4 address. I am not sure if
IPv6 is supported by the protocol. In
theory it could be, but this CLI only
supports IPv4 for now.

Please also note that this operation
does not interact with the switches.
Therefore it will ignore the retries
command line flag.`,
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
		fmt.Fprintf(w, "INTERFACE\tMAC\tIP\n")

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

			// Skip if interface has no valid IPv4.
			ip, err := nsdp.GetInterfaceIPv4(&iface)
			if err != nil || ip == nil {
				continue
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", iface.Name, iface.HardwareAddr.String(), ip.String())
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
