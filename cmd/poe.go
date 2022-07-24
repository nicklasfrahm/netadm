package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/nicklasfrahm/netadm/pkg/nsdp"
	"github.com/spf13/cobra"
)

var poeCmd = &cobra.Command{
	Use:   "poe",
	Short: "Change power over Ethernet settings",
	Long: `Change the power over Ethernet settings
of a network device.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	SilenceUsage: true,
}

var poeOnCmd = &cobra.Command{
	Use:   "on <device> [ports]",
	Short: "Enable power over Ethernet",
	Long: `Enable power over Ethernet for a single or multiple ports.

You may specify the ports as a comma-separated list.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
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
			request := nsdp.NewMessage(nsdp.WriteRequest)

			// Add request records.
			name := fmt.Sprintf("test-%d", time.Now().Unix()%1e3)

			// Obscure passwork using "secret".
			secret := "NtgrSmartSwitchRock"
			passwordBytes := []byte(secret)
			obscured := make([]byte, len(secret))
			for i := 0; i < len(obscured); i++ {
				if i < len(passwordBytes) {
					obscured[i] = passwordBytes[i] ^ secret[i%len(secret)]
				} else {
					obscured[i] = 0
				}
			}

			fmt.Println("obscured:", obscured)

			request.Records = append(request.Records, nsdp.Record{
				ID:    nsdp.RecordName.ID,
				Len:   uint16(len(name)),
				Value: []byte(name),
			})

			request.Records = append(request.Records, nsdp.Record{
				ID:    nsdp.RecordPassword.ID,
				Len:   uint16(len(password)),
				Value: []byte(password),
			})

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

		// TODO: Should we print a result?

		return nil
	},
}

var poeOffCmd = &cobra.Command{
	Use:   "off <device> [ports]",
	Short: "Disable power over Ethernet",
	Long: `Disable power over Ethernet for a single or multiple ports.

You may specify the ports as a comma-separated list.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func init() {
	poeOnCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	poeOnCmd.MarkFlagRequired("interface")
	poeOffCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	poeOffCmd.MarkFlagRequired("interface")

	poeCmd.AddCommand(poeOnCmd)
	poeCmd.AddCommand(poeOffCmd)

	rootCmd.AddCommand(poeCmd)
}
