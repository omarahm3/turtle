package cmd

import (
	"fmt"
	"os"

	"github.com/omarahm3/turtle/pkg/helpers"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var (
	nethogs   bool
	bandwhich bool
	db        *buntdb.DB
	dbPath    = "./turtle.db"
	rootCmd   = &cobra.Command{
		Use:   "turtle",
		Short: "Log nethogs traffic per processes and applications",
		Run:   runner,
	}
)

func Init() {
	if !helpers.IsRoot() {
		fmt.Println("you need to run this with sudo")
		os.Exit(1)
	}

	_, err := helpers.SnifferExists()
	check(err)
	db, err = buntdb.Open(dbPath)
	check(err)
	db.CreateIndex("apps", "*:app", buntdb.IndexString)
	db.CreateIndex("processes", "*:process", buntdb.IndexString)

	rootCmd.Flags().BoolVarP(&nethogs, "nethogs", "n", true, "use nethogs as a underlayer processor")
	rootCmd.Flags().BoolVarP(&bandwhich, "bandwhich", "b", false, "use bandwhich as a underlayer processor")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check(err error) {
	if err == nil {
		return
	}

	fmt.Printf("error occurred: %s\n", err.Error())
	os.Exit(1)
}

func fatalPrint(s string) {
	fmt.Println(s)
	os.Exit(1)
}
