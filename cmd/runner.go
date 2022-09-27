package cmd

import (
	"errors"
	"fmt"

	"github.com/omarahm3/turtle/pkg/sniffer"
	"github.com/omarahm3/turtle/pkg/sniffer/processors"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var processor processors.Processor

func runner(cmd *cobra.Command, args []string) {
	if bandwhich {
		processor = processors.New(processors.TYPE_BANDWHICH)
	} else {
		processor = processors.New(processors.TYPE_NETHOGS)
	}

	if !processor.CanRunCommand() {
		check(fmt.Errorf("command %q does not exist on the system, please make sure it is installed properly", processor.GetCommand().Args[0]))
	}

	sl := make(chan sniffer.SniffLog)
	go sniffer.Sniff(sl, processor)

	for {
		update(<-sl)
	}
}

func update(sl sniffer.SniffLog) {
	var appValue string
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		appValue, err = tx.Get(fmt.Sprintf("%s:app", sl.AppHash))
		if err != nil && errors.Is(err, buntdb.ErrNotFound) {
			return nil
		}
		return err
	})
	check(err)

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprintf("%s:process", sl.TotalPathHash), sl.String(), nil)
		if err != nil {
			return err
		}

		// In case app value does not exist
		if appValue == "" {
			_, _, err := tx.Set(fmt.Sprintf("%s:app", sl.AppHash), sl.NewAppLog().String(), nil)
			return err
		}

		// if it exists then just update the metrics
		oldAppLog := sniffer.ToAppLog(appValue)
		newAppLog := sl.NewAppLog()
		newAppLog.Sent += oldAppLog.Sent
		newAppLog.Received += oldAppLog.Received
		_, _, err = tx.Set(fmt.Sprintf("%s:app", newAppLog.AppHash), newAppLog.String(), nil)

		return err
	})
	check(err)
}
