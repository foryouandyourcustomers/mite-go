package cmd

import (
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func init() {
	projectsCommand.AddCommand(listProjectsCommand)
	rootCmd.AddCommand(projectsCommand)
}

var projectsCommand = &cobra.Command{
	Use:   "projects",
	Short: "list & adds projects",
	RunE:  listProjectsCommand.RunE,
}

var listProjectsCommand = &cobra.Command{
	Use:   "list",
	Short: "list projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := deps.miteApi.Projects()
		if err != nil {
			return err
		}

		t := tabby.New()
		t.AddHeader("id", "name", "notes")
		for _, project := range projects {
			t.AddLine(project.Id, project.Name, project.Note)
		}
		t.Print()
		return nil
	},
}
