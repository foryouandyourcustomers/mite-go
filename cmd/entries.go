package cmd

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
	listTo          string
	listFrom        string
	listOrder       string
	createDate      string
	createDuration  time.Duration
	createNote      string
	createProjectId string
	createServiceId string
)

func init() {
	now := time.Now()
	defaultFrom := now.AddDate(0, 0, -7)
	defaultDuration, err := time.ParseDuration("0m")
	if err != nil {
		panic(err)
	}
	// list
	entriesListCommand.Flags().StringVarP(&listTo, "to", "t", now.Format("2006-01-02"), "list only entries until date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listFrom, "from", "f", defaultFrom.Format("2006-01-02"), "list only entries starting at date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listOrder, "order", "o", "asc", "list only entries starting at date (in YYYY-MM-DD format)")
	entriesCommand.AddCommand(entriesListCommand)
	// flags for create
	entriesCreateCommand.Flags().StringVarP(&createDate, "date", "D", now.Format("2006-01-02"), "day for which to create entry (in YYYY-MM-DD format)")
	entriesCreateCommand.Flags().DurationVarP(&createDuration, "duration", "d", defaultDuration, "duration of entry (format examples: '1h15m' or '300m' or '6h')")
	entriesCreateCommand.Flags().StringVarP(&createNote, "note", "n", "", "a note describing what was worked on")
	entriesCreateCommand.Flags().StringVarP(&createProjectId, "projectid", "p", "", "project id for time entry (HINT: use the 'project' sub-command to find the id)")
	entriesCreateCommand.Flags().StringVarP(&createServiceId, "serviceid", "s", "", "service id for time entry (HINT: use the 'service' sub-command to find the id)")
	entriesCommand.AddCommand(entriesCreateCommand)
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

		entries, err := deps.miteApi.TimeEntries(&mite.TimeEntryQuery{
			To:        &to,
			From:      &from,
			Direction: direction,
		})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		printEntries(entries)
	},
}

func printEntries(entries []*mite.TimeEntry) {
	t := tabby.New()
	t.AddHeader("id", "notes", "date", "time", "project,service")
	for _, entry := range entries {
		trimmedNotes := strings.Replace(entry.Note, "\r\n", ",", -1)
		shortenedNotes := fmt.Sprintf("%.50s", trimmedNotes)
		shortenedProjectService := fmt.Sprintf("%.50s", entry.ProjectName+","+entry.ServiceName)
		t.AddLine(entry.Id, shortenedNotes, entry.Date, entry.Duration.String(), shortenedProjectService)
	}
	t.Print()
}

var entriesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "create time entries",
	Run: func(cmd *cobra.Command, args []string) {
		if createProjectId == "" || createServiceId == "" {
			_, _ = fmt.Fprintln(os.Stderr, "please set both the project AND service id")
			return
		}

		cDate, err := time.Parse("2006-01-02", createDate)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		timeEntry := mite.TimeEntryCommand{
			Date:      &cDate,
			Duration:  &createDuration,
			Note:      createNote,
			ProjectId: createProjectId,
			ServiceId: createServiceId,
		}

		entry, err := deps.miteApi.CreateTimeEntry(&timeEntry)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		printEntries([]*mite.TimeEntry{entry})
	},
}
