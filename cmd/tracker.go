package cmd

import (
	"errors"
	"fmt"
	"github.com/leanovate/mite-go/date"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
)

var (
	trackerTimeEntryId string
)

func init() {
	trackerCommand.AddCommand(trackerStatusCommand)
	trackerStartCommand.Flags().StringVarP(&trackerTimeEntryId, "id", "i", "", "the time entry id to (re)start a tracker for (default: latest time entry for today)")
	trackerCommand.AddCommand(trackerStartCommand)
	trackerStopCommand.Flags().StringVarP(&trackerTimeEntryId, "id", "i", "", "the time entry id to stop a tracker for (default: latest time entry for today)")
	trackerCommand.AddCommand(trackerStopCommand)
	rootCmd.AddCommand(trackerCommand)
}

var trackerCommand = &cobra.Command{
	Use:   "tracker",
	Short: "starts, stops and shows the status of the tracker",
	RunE:  trackerStatusCommand.RunE,
}

var trackerStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "shows the status of the time tracker",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		panic("implement trackerStatusCommand")
		return nil
	},
}

var trackerStartCommand = &cobra.Command{
	Use:   "start",
	Short: "starts the time tracker for a time entry",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if trackerTimeEntryId == "" {
			trackerTimeEntryId, err = fetchLatestTimeEntryForToday()
			if err != nil {
				return err
			}
		}
		fmt.Printf("passed id: %s\n", trackerTimeEntryId)
		panic("implement trackerStartCommand")
		return nil
	},
}

func fetchLatestTimeEntryForToday() (string, error) {
	today := date.Today()

	entries, err := deps.miteApi.TimeEntries(&mite.TimeEntryQuery{
		To:        &today,
		From:      &today,
		Direction: "desc",
	})
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", errors.New("no time entries for today found")
	}

	return entries[0].Id, nil
}

var trackerStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "stops the time tracker for a time entry",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if trackerTimeEntryId == "" {
			trackerTimeEntryId, err = fetchLatestTimeEntryForToday()
			if err != nil {
				return err
			}
		}
		fmt.Printf("passed id: %s\n", trackerTimeEntryId)
		panic("implement trackerStopCommand")
		return nil
	},
}
