package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Read configuration keys",
	Long: `A command that allows you to read the
list of specified configuration keys.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("get called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
