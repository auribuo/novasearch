package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the catalog cache",
	Long:  `Purge the catalog cache`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Purging all catalogs")
		} else {
			for _, arg := range args {
				fmt.Printf("Purging catalog %s", arg)
			}
		}
	},
	Args: cobra.MaximumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ugc", "ned"}, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
