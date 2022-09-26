package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

const (
	apps_type      = "apps"
	processes_type = "processes"
)

var (
	apps        bool
	processes   bool
	listCommand = &cobra.Command{
		Use:   "list",
		Short: "List collected information",
		Run:   list,
	}
)

func list(cmd *cobra.Command, args []string) {
	err := db.View(func(tx *buntdb.Tx) error {
		cType := apps_type
		if processes {
			cType = processes_type
		}

		return tx.Ascend(cType, func(key, value string) bool {
			fmt.Println(value)
			return true
		})
	})
	check(err)
}

func init() {
	listCommand.PersistentFlags().BoolVarP(&apps, "apps", "a", true, "List all apps usage")
	listCommand.PersistentFlags().BoolVarP(&processes, "processes", "p", false, "List all processes usage")
	rootCmd.AddCommand(listCommand)
}
