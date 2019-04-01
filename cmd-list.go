package main

import (
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func init() {
	listCommand.AddCommand(listProjectsCommand)
	listCommand.AddCommand(listServicesCommand)
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
