package cmd

import (
	"errors"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/leanovate/mite-go/date"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
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
	editTimeEntryId string
	editDate        string
	editDuration    string
	editNote        string
	editProjectId   string
	editServiceId   string
)

func init() {
	today := date.Today()
	defaultFrom := today.Add(0, 0, -7)
	defaultDuration, err := time.ParseDuration("0m")
	if err != nil {
		panic(err)
	}
	// list
	entriesListCommand.Flags().StringVarP(&listTo, "to", "t", today.String(), "list only entries until date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listFrom, "from", "f", defaultFrom.String(), "list only entries starting at date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listOrder, "order", "o", "asc", "list only entries starting at date (in YYYY-MM-DD format)")
	entriesCommand.AddCommand(entriesListCommand)
	// create
	entriesCreateCommand.Flags().StringVarP(&createDate, "date", "D", today.String(), "day for which to create entry (in YYYY-MM-DD format)")
	entriesCreateCommand.Flags().DurationVarP(&createDuration, "duration", "d", defaultDuration, "duration of entry (format examples: '1h15m' or '300m' or '6h')")
	entriesCreateCommand.Flags().StringVarP(&createNote, "note", "n", "", "a note describing what was worked on")
	entriesCreateCommand.Flags().StringVarP(&createProjectId, "projectid", "p", "", "project id for time entry (HINT: use the 'project' sub-command to find the id)")
	entriesCreateCommand.Flags().StringVarP(&createServiceId, "serviceid", "s", "", "service id for time entry (HINT: use the 'service' sub-command to find the id)")
	entriesCommand.AddCommand(entriesCreateCommand)
	// edit
	entriesEditCommand.Flags().StringVarP(&editDate, "date", "D", now.Format("2006-01-02"), "day for which to edit entry (in YYYY-MM-DD format)")
	entriesEditCommand.Flags().StringVarP(&editDuration, "duration", "d", "", "duration of entry (format examples: '1h15m' or '300m' or '6h')")
	entriesEditCommand.Flags().StringVarP(&editNote, "note", "n", "", "a note describing what was worked on")
	entriesEditCommand.Flags().StringVarP(&editTimeEntryId, "id", "i", "", "the time entry id to edit")
	entriesEditCommand.Flags().StringVarP(&editProjectId, "projectid", "p", "", "project id for time entry (HINT: use the 'project' sub-command to find the id)")
	entriesEditCommand.Flags().StringVarP(&editServiceId, "serviceid", "s", "", "service id for time entry (HINT: use the 'service' sub-command to find the id)")
	entriesCommand.AddCommand(entriesEditCommand)
	rootCmd.AddCommand(entriesCommand)
}

var entriesCommand = &cobra.Command{
	Use:   "entries",
	Short: "lists & adds time entries",
	RunE:  entriesListCommand.RunE,
}

var entriesListCommand = &cobra.Command{
	Use:   "list",
	Short: "list time entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		direction := listOrder

		to, err := date.Parse(listTo)
		if err != nil {
			return err
		}
		from, err := date.Parse(listFrom)
		if err != nil {
			return err
		}

		entries, err := deps.miteApi.TimeEntries(&mite.TimeEntryQuery{
			To:        &to,
			From:      &from,
			Direction: direction,
		})
		if err != nil {
			return err
		}

		printEntries(entries)
		return nil
	},
}

func printEntries(entries []*mite.TimeEntry) {
	t := tabby.New()
	t.AddHeader("id", "notes", "date", "time", "project", "service")
	for _, entry := range entries {
		trimmedNotes := strings.Replace(entry.Note, "\r\n", ",", -1)
		shortenedNotes := fmt.Sprintf("%.50s", trimmedNotes)
		shortenedProject := fmt.Sprintf("%.25s", entry.ProjectName)
		shortenedService := fmt.Sprintf("%.25s", entry.ServiceName)
		t.AddLine(entry.Id, shortenedNotes, entry.Date, entry.Duration.String(), shortenedProject, shortenedService)
	}
	t.Print()
}

var entriesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "creates a time entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if createProjectId == "" {
			createProjectId = deps.conf.Get("projectId")
		}

		if createServiceId == "" {
			createServiceId = deps.conf.Get("serviceId")
		}

		if createProjectId == "" || createServiceId == "" {
			return errors.New("please set both the project AND service id (either via arguments or config)")
		}

		cDate, err := date.Parse(createDate)
		if err != nil {
			return err
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
			return err
		}

		printEntries([]*mite.TimeEntry{entry})
		return nil
	},
}

var entriesEditCommand = &cobra.Command{
	Use:   "edit",
	Short: "edits a time entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := deps.miteApi.TimeEntry(editTimeEntryId)
		if err != nil {
			return err
		}

		// use retrieved values as defaults
		timeEntry := mite.TimeEntryCommand{
			Date:      &entry.Date,
			Duration:  &entry.Duration,
			Note:      entry.Note,
			ProjectId: entry.ProjectId,
			ServiceId: entry.ServiceId,
		}

		// override only fields affected by set parameters of edit
		if editDate != "" {
			eDate, err := time.Parse("2006-01-02", editDate)
			if err != nil {
				return err
			}
			timeEntry.Date = &eDate
		}

		if editDuration != "" {
			eDuration, err := time.ParseDuration(editDuration)
			if err != nil {
				return err
			}
			timeEntry.Duration = &eDuration
		}

		if editNote != "" {
			timeEntry.Note = editNote
		}

		if editProjectId != "" {
			timeEntry.ProjectId = editProjectId
		}

		if editServiceId != "" {
			timeEntry.ServiceId = editServiceId
		}

		err = deps.miteApi.EditTimeEntry(editTimeEntryId, &timeEntry)
		if err != nil {
			return err
		}

		entry, err = deps.miteApi.TimeEntry(editTimeEntryId)
		if err != nil {
			return err
		}

		printEntries([]*mite.TimeEntry{entry})
		return nil
	},
}
