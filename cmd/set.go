package cmd

import (
	"errors"
	"fmt"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [key=value ...]",
	Short: "Write configuration keys",
	Long: `A command that allows you to write the
list of specified configuration keys.

You may run the "keys" subcommand
to see a list of available keys.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if all keys are valid.
		for _, key := range args {
			if nsdp.RecordTypeNames[key] == nil {
				return fmt.Errorf("unknown configuration key %s", key)
			}
		}

		// TODO: Implement "set" command logic.
		return errors.New("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
