package main

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

func init() {
	entriesCommand.AddCommand(entriesListCommand)
	rootCmd.AddCommand(entriesCommand)
}

var entriesCommand = &cobra.Command{
	Use:   "entries",
	Short: "lists & adds time entries",
	Run:   entriesListCommand.Run,
}

var entriesListCommand = &cobra.Command{
	Use:   "list",
	Short: "list time entries",
	Run: func(cmd *cobra.Command, args []string) {
		api := mite.NewMiteApi(configGetApiUrl(), configGetApiKey())
		to := time.Now()
		from := to.AddDate(0, 0, -7)
		direction := mite.DirectionAsc

		entries, err := api.TimeEntries(&mite.TimeEntryParameters{
			To:        &to,
			From:      &from,
			Direction: &direction,
		})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		t := tabby.New()
		t.AddHeader("id", "notes", "date", "time", "project,service")
		for _, entry := range entries {
			trimmedNotes := strings.Replace(entry.Note, "\r\n", ",", -1)
			shortendNotes := fmt.Sprintf("%.50s", trimmedNotes)
			shortenedProjectService := fmt.Sprintf("%.50s", entry.ProjectName+","+entry.ServiceName)
			t.AddLine(entry.Id, shortendNotes, entry.Date, entry.Duration.String(), shortenedProjectService)
		}
		t.Print()
	},
}
