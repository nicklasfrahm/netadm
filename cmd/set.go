package cmd

import (
	"errors"
	"fmt"

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
		}

		id := args[0]
		if id == "all" {
			return errors.New("writing to all devices is not supported")
		}

		_ = args[1:]

		devices, err := nsdp.Get(id, []string{"mac", "encryptionmode"}, opts...)
		if err != nil {
			return err
		}
		encryptionMode := devices[0].EncryptionMode

		nonce := make([]byte, 4)
		if encryptionMode == nsdp.EncryptionModeHash32 || encryptionMode == nsdp.EncryptionModeHash64 {
			devs, err := nsdp.Get(id, []string{"mac", "encryptionnonce"}, opts...)
			if err != nil {
				return err
			}

			if len(nonce) == 0 {
				return errors.New("failed to fetch encryption nonce")
			}

			copy(nonce, devs[0].EncryptionNonce)
		}

		encryptedPassword, err := nsdp.EncryptPassword(encryptionMode, devices[0].MAC, nonce, []byte(password))
		if err != nil {
			return err
		}

		fmt.Println(encryptedPassword)

		// Print results.
		// fmt.Table(os.Stdout, devices, []string{"mac", "encryptionmode"})

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
