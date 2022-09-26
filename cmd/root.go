package cmd

import (
	"fmt"
	"os"

	"github.com/omarahm3/turtle/pkg/helpers"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var (
	db      *buntdb.DB
	dbPath  = "./turtle.db"
	rootCmd = &cobra.Command{
		Use:   "turtle",
		Short: "Log nethogs traffic per processes and applications",
		Run:   runner,
	}
)

func Init() {
	_, err := helpers.SnifferExists()
	check(err)
	db, err = buntdb.Open(dbPath)
	check(err)
	db.CreateIndex("apps", "*:app", buntdb.IndexString)
	db.CreateIndex("processes", "*:process", buntdb.IndexString)

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
