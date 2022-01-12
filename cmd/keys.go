package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
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
		// Sort keys by RecordTypeID to get consistent results.
		keyIDs := make([]int, 0, len(nsdp.RecordTypeNames))
		for _, recordType := range nsdp.RecordTypeNames {
			keyIDs = append(keyIDs, int(recordType.ID))
		}
		sort.Ints(keyIDs)

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "ID\tName\tExample\n")

		// Print a list of all available configuration keys.
		for _, keyID := range keyIDs {
			recordType := nsdp.RecordTypeIDs[nsdp.RecordTypeID(keyID)]
			fmt.Fprintf(w, "0x%04X\t%s\t%v\n", recordType.ID, strings.ToLower(recordType.Name), recordType.Example)
		}

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
