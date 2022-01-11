package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicklasfrahm/nsdp/pkg/nsdp"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List available configuration keys",
	Long: `A command that allows you to list all
available configuration keys.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Sort keys by record ID?

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "ID\tName\tExample\n")

		// Print a list of all available configuration keys.
		for _, key := range nsdp.RecordNames {
			fmt.Fprintf(w, "0x%04X\t%s\t%s\n", key.ID, key.Name, key.Example)
		}

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
