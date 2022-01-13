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
		ids := make([]int, 0, len(nsdp.RecordTypeByName))
		for _, rt := range nsdp.RecordTypeByName {
			ids = append(ids, int(rt.ID))
		}
		sort.Ints(ids)

		// Create table with tabwriter.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "ID\tNAME\tEXAMPLE\n")

		// Print a list of all available configuration keys.
		for _, id := range ids {
			// Fetch record type by ID.
			rt := nsdp.RecordTypeByID[nsdp.RecordTypeID(id)]
			fmt.Fprintf(w, "0x%04X\t%s\t%v\n", rt.ID, strings.ToLower(rt.Name), rt.Example)
		}

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
