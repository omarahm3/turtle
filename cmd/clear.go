package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	force        bool
	clearCommand = &cobra.Command{
		Use:   "clear",
		Short: "Clear database",
		Run:   clear,
	}
)

func clear(cmd *cobra.Command, args []string) {
	if !force {
		fmt.Println("You need to use '-f|--force' if you're sure")
		return
	}

	err := os.Remove(getDbPath())
	check(err)
}

func init() {
	clearCommand.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force")
	rootCmd.AddCommand(clearCommand)
}
