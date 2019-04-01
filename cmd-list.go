package main

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"strings"
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
		projects := apiGetProjects()
		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, project := range projects {
			p := project.ProjectBody
			t.AddLine(p.Id, p.Name, p.Note)
		}
		t.Print()
	},
}

var listServicesCommand = &cobra.Command{
	Use:   "services",
	Short: "list services",
	Run: func(cmd *cobra.Command, args []string) {
		services := apiGetServices()
		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, project := range services {
			s := project.ServiceBody
			t.AddLine(s.Id, s.Name, s.Note)
		}
		t.Print()
	},
}

var listTimeEntriesCommand = &cobra.Command{
	Use:   "entries",
	Short: "list time entries",
	Run: func(cmd *cobra.Command, args []string) {
		entries := apiGetEntries()
		t := tabby.New()
		t.AddHeader("id", "notes", "date", "time", "project,service")
		for _, entry := range entries {
			s := entry.TimeEntryBody
			trimmedNotes := strings.Replace(s.Note, "\r\n", ",", -1)
			shortendNotes := fmt.Sprintf("%.50s", trimmedNotes)
			shortenedProjectService := fmt.Sprintf("%.50s", s.ProjectName+","+s.ServiceName)
			t.AddLine(s.Id, shortendNotes, s.Date, formatMinutesToHuman(s.Minutes), shortenedProjectService)
		}
		t.Print()
	},
}

func formatMinutesToHuman(minutes int) string {
	if minutes > 60 {
		hours := minutes / 60
		return fmt.Sprintf("%.2dh:%.2dm", hours, minutes-hours*60)
	}

	return fmt.Sprintf("0h:%.2dm", minutes)
}
