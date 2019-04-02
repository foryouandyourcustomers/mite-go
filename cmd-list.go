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
	listCommand.AddCommand(listProjectsCommand)
	listCommand.AddCommand(listServicesCommand)
	listCommand.AddCommand(listTimeEntriesCommand)
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "list entries, projects and roles",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: list entries for last 7 days by default
	},
}

var listProjectsCommand = &cobra.Command{
	Use:   "projects",
	Short: "list projects",
	Run: func(cmd *cobra.Command, args []string) {
		api := mite.NewMiteApi(configGetApiUrl(), configGetApiKey())
		projects, err := api.Projects()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, project := range projects {
			t.AddLine(project.Id, project.Name, project.Note)
		}
		t.Print()
	},
}

var listServicesCommand = &cobra.Command{
	Use:   "services",
	Short: "list services",
	Run: func(cmd *cobra.Command, args []string) {
		api := mite.NewMiteApi(configGetApiUrl(), configGetApiKey())
		services, err := api.Services()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, service := range services {
			t.AddLine(service.Id, service.Name, service.Note)
		}
		t.Print()
	},
}

var listTimeEntriesCommand = &cobra.Command{
	Use:   "entries",
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
