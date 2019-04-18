package cmd

import (
	"errors"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/leanovate/mite-go/domain"
	"github.com/spf13/cobra"
	"strings"
)

var (
	listTo            string
	listFrom          string
	listOrder         string
	createDate        string
	createMinutes     string
	createNote        string
	createProjectId   string
	createServiceId   string
	createActivity    string
	editTimeEntryId   string
	editDate          string
	editMinutes       string
	editNote          string
	editProjectId     string
	editServiceId     string
	editActivity      string
	deleteTimeEntryId string
)

func init() {
	today := domain.Today()
	defaultFrom := today.Add(0, 0, -7)
	defaultMinutes := domain.NewMinutes(0).String()

	// list
	entriesListCommand.Flags().StringVarP(&listTo, "to", "t", today.String(), "list only entries until date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listFrom, "from", "f", defaultFrom.String(), "list only entries starting at date (in YYYY-MM-DD format)")
	entriesListCommand.Flags().StringVarP(&listOrder, "order", "o", "asc", "list only entries starting at date (in YYYY-MM-DD format)")
	entriesCommand.AddCommand(entriesListCommand)
	// create
	entriesCreateCommand.Flags().StringVarP(&createDate, "date", "D", today.String(), "day for which to create entry (in YYYY-MM-DD format)")
	entriesCreateCommand.Flags().StringVarP(&createMinutes, "duration", "d", defaultMinutes, "duration of entry (format examples: '1h15m' or '300m' or '6h')")
	entriesCreateCommand.Flags().StringVarP(&createNote, "note", "n", "", "a note describing what was worked on")
	entriesCreateCommand.Flags().StringVarP(&createProjectId, "projectid", "p", "", "project id for time entry (HINT: use the 'project' sub-command to find the id)")
	entriesCreateCommand.Flags().StringVarP(&createServiceId, "serviceid", "s", "", "service id for time entry (HINT: use the 'service' sub-command to find the id)")
	entriesCreateCommand.Flags().StringVarP(&createActivity, "activity", "a", "", "activity describing a specific project and service combination")
	entriesCommand.AddCommand(entriesCreateCommand)
	// edit
	entriesEditCommand.Flags().StringVarP(&editDate, "date", "D", "", "day for which to edit entry (in YYYY-MM-DD format)")
	entriesEditCommand.Flags().StringVarP(&editMinutes, "duration", "d", "", "duration of entry (format examples: '1h15m' or '300m' or '6h')")
	entriesEditCommand.Flags().StringVarP(&editNote, "note", "n", "", "a note describing what was worked on")
	entriesEditCommand.Flags().StringVarP(&editTimeEntryId, "id", "i", "", "the time entry id to edit")
	entriesEditCommand.Flags().StringVarP(&editProjectId, "projectid", "p", "", "project id for time entry (HINT: use the 'project' sub-command to find the id)")
	entriesEditCommand.Flags().StringVarP(&editServiceId, "serviceid", "s", "", "service id for time entry (HINT: use the 'service' sub-command to find the id)")
	entriesEditCommand.Flags().StringVarP(&editActivity, "activity", "a", "", "activity describing a specific project and service combination")
	entriesCommand.AddCommand(entriesEditCommand)
	// delete
	entriesDeleteCommand.Flags().StringVarP(&deleteTimeEntryId, "id", "i", "", "the time entry id to delete")
	entriesCommand.AddCommand(entriesDeleteCommand)
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

		to, err := domain.ParseLocalDate(listTo)
		if err != nil {
			return err
		}
		from, err := domain.ParseLocalDate(listFrom)
		if err != nil {
			return err
		}

		entries, err := application.MiteApi.TimeEntries(&domain.TimeEntryQuery{
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

func printEntries(entries []*domain.TimeEntry) {
	t := tabby.New()
	t.AddHeader("id", "notes", "date", "time", "project", "service")
	for _, entry := range entries {
		trimmedNotes := strings.Replace(entry.Note, "\r\n", ",", -1)
		shortenedNotes := fmt.Sprintf("%.50s", trimmedNotes)
		shortenedProject := fmt.Sprintf("%.25s", entry.ProjectName)
		shortenedService := fmt.Sprintf("%.25s", entry.ServiceName)
		t.AddLine(entry.Id, shortenedNotes, entry.Date, entry.Minutes.String(), shortenedProject, shortenedService)
	}
	t.Print()
}

var entriesCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "creates a time entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectId, serviceId := projectAndServiceId()
		if projectId == "" || serviceId == "" {
			return errors.New("please set both the project AND service id (either via arguments or config)")
		}

		cDate, err := domain.ParseLocalDate(createDate)
		if err != nil {
			return err
		}
		cMinutes, err := domain.ParseMinutes(createMinutes)
		if err != nil {
			return err
		}
		cProjectId, err := domain.ParseProjectId(projectId)
		if err != nil {
			return err
		}
		cServiceId, err := domain.ParseServiceId(serviceId)
		if err != nil {
			return err
		}

		timeEntry := domain.TimeEntryCommand{
			Date:      &cDate,
			Minutes:   &cMinutes,
			Note:      createNote,
			ProjectId: cProjectId,
			ServiceId: cServiceId,
		}

		entry, err := application.MiteApi.CreateTimeEntry(&timeEntry)
		if err != nil {
			return err
		}

		printEntries([]*domain.TimeEntry{entry})
		return nil
	},
}

func projectAndServiceId() (projectId, servicesId string) {
	if createProjectId == "" && createActivity != "" {
		activity := application.Conf.GetActivity(createActivity)
		createProjectId = activity.ProjectId
	}

	if createServiceId == "" && createActivity != "" {
		activity := application.Conf.GetActivity(createActivity)
		createServiceId = activity.ServiceId
	}

	if createProjectId == "" {
		createProjectId = application.Conf.Get("projectId")
	}

	if createServiceId == "" {
		createServiceId = application.Conf.Get("serviceId")
	}

	return createProjectId, createServiceId
}

var entriesEditCommand = &cobra.Command{
	Use:   "edit",
	Short: "edits a time entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		entryId, err := domain.ParseTimeEntryId(editTimeEntryId)
		if err != nil {
			return err
		}

		entry, err := application.MiteApi.TimeEntry(entryId)
		if err != nil {
			return err
		}

		// use retrieved values as defaults
		command := domain.TimeEntryCommand{
			Date:      &entry.Date,
			Minutes:   &entry.Minutes,
			Note:      entry.Note,
			ProjectId: entry.ProjectId,
			ServiceId: entry.ServiceId,
		}

		// override only fields affected by set parameters of edit
		if editDate != "" {
			eDate, err := domain.ParseLocalDate(editDate)
			if err != nil {
				return err
			}
			command.Date = &eDate
		}

		if editMinutes != "" {
			eMinutes, err := domain.ParseMinutes(editMinutes)
			if err != nil {
				return err
			}
			command.Minutes = &eMinutes
		}

		if editNote != "" {
			command.Note = editNote
		}

		if editActivity != "" {
			activity := application.Conf.GetActivity(editActivity)

			projectId, err := domain.ParseProjectId(activity.ProjectId)
			if err != nil {
				return err
			}
			command.ProjectId = projectId

			serviceId, err := domain.ParseServiceId(activity.ServiceId)
			if err != nil {
				return err
			}
			command.ServiceId = serviceId
		}

		if editProjectId != "" && command.ProjectId == 0 {
			projectId, err := domain.ParseProjectId(editProjectId)
			if err != nil {
				return err
			}
			command.ProjectId = projectId
		}

		if editServiceId != "" && command.ProjectId == 0 {
			serviceId, err := domain.ParseServiceId(editServiceId)
			if err != nil {
				return err
			}
			command.ServiceId = serviceId
		}

		err = application.MiteApi.EditTimeEntry(entryId, &command)
		if err != nil {
			return err
		}

		entry, err = application.MiteApi.TimeEntry(entryId)
		if err != nil {
			return err
		}

		printEntries([]*domain.TimeEntry{entry})
		return nil
	},
}

var entriesDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "deletes a time entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		entryId, err := domain.ParseTimeEntryId(deleteTimeEntryId)
		if err != nil {
			return err
		}

		return application.MiteApi.DeleteTimeEntry(entryId)
	},
}
