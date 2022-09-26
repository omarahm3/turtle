package cmd

import (
	"fmt"
	"math"
	"strings"

	"github.com/omarahm3/turtle/pkg/sniffer"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

const (
	apps_type        = "apps"
	processes_type   = "processes"
	max_chars_length = 40
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
			if cType == apps_type {
				printApp(value)
				return true
			}
			printProcess(value)
			return true
		})
	})
	check(err)
}

func printApp(v string) {
	a := sniffer.ToAppLog(v)
	fmt.Printf("%s\t%fMB\t%fMB\n", trim(a.App), a.Sent, a.Received)
}

func printProcess(v string) {
	s := sniffer.ToSniffLog(v)
	fmt.Printf("%s\t%fMB\t%fMB\n", trim(s.App), s.Sent, s.Received)
}

func trim(s string) string {
	min := int(math.Min(float64(len(s)), max_chars_length))
	if len(s) < max_chars_length {
		return fmt.Sprintf("%s%s", s, strings.Repeat(" ", max_chars_length-len(s)))
	}
	return fmt.Sprintf("%s...", s[0:min])
}

func init() {
	listCommand.PersistentFlags().BoolVarP(&apps, "apps", "a", true, "List all apps usage")
	listCommand.PersistentFlags().BoolVarP(&processes, "processes", "p", false, "List all processes usage")
	rootCmd.AddCommand(listCommand)
}
