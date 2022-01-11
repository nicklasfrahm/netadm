package cmd

import (
	"errors"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [keys]",
	Short: "Read configuration keys",
	Long: `A command that allows you to read the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if all keys are valid.
		for _, key := range args {
			if nsdp.RecordNames[key] == nil {
				return errors.New("unknown configuration key")
			}
		}

		// TODO: Implement get command logic.

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}