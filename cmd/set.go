package cmd

import (
	"errors"
	"strings"

	"github.com/nicklasfrahm/netadm/pkg/nsdp"
	"github.com/spf13/cobra"
)

var password string

var setCmd = &cobra.Command{
	Use:   "set <device> [key=value ...]",
	Short: "Write configuration keys",
	Long: `A command that allows you to write the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := []nsdp.Option{
			nsdp.WithInterfaceName(interfaceName),
			nsdp.WithRetries(retries),
			nsdp.WithTimeout(timeout),
			nsdp.WithPassword(password),
		}

		id := args[0]
		if id == "all" {
			return errors.New("writing to all devices is not supported")
		}

		// Separate the key-value pairs into a map.
		values := make(map[string]string)
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				values[strings.ToLower(parts[0])] = parts[1]
			}
		}

		_, err := nsdp.Set(id, values, opts...)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "name of the interface to use")
	setCmd.MarkFlagRequired("interface")
	setCmd.Flags().StringVarP(&password, "password", "p", "", "password to use for authentication")
	setCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(setCmd)
}
