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

var (
	listTo    string
	listFrom  string
	listOrder string
)

func init() {
	defaultTo := time.Now()
	defaultFrom := defaultTo.AddDate(0, 0, -7)
	entriesListCommand.Flags().StringVarP(&listTo, "to", "t", defaultTo.Format("2006-01-02"), "list only entries until date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listFrom, "from", "f", defaultFrom.Format("2006-01-02"), "list only entries starting at date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listOrder, "order", "o", "asc", "list only entries starting at date (in YYYY-MM-DD format)")
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

		direction := listOrder

		to, err := time.Parse("2006-01-02", listTo)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}
		from, err := time.Parse("2006-01-02", listFrom)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		entries, err := api.TimeEntries(&mite.TimeEntryParameters{
			To:        &to,
			From:      &from,
			Direction: direction,
		})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		t := tabby.New()
		t.AddHeader("id", "notes", "date", "time", "project,service")
		for _, entry := range entries {
			trimmedNotes := strings.Replace(entry.Note, "\r\n", ",", -1)
			shortenedNotes := fmt.Sprintf("%.50s", trimmedNotes)
			shortenedProjectService := fmt.Sprintf("%.50s", entry.ProjectName+","+entry.ServiceName)
			t.AddLine(entry.Id, shortenedNotes, entry.Date, entry.Duration.String(), shortenedProjectService)
		}
		t.Print()
	},
}
