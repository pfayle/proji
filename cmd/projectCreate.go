package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/models"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create LABEL NAME [NAME...]",
	Short: "Create one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("at least one label and name have to be given")
		}
		label := args[0]
		projects := args[1:]

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Load class once for all projects
		class, err := projiEnv.Svc.LoadClass(label)
		if err != nil {
			return err
		}

		for _, name := range projects {
			fmt.Printf("\n> Creating project %s\n", name)

			err := createProject(name, cwd, projiEnv.ConfigFolderPath, class)
			if err != nil {
				fmt.Printf(" -> Failed: %v\n", err)

				if err.Error() == "Project already exists" {
					if !helper.WantTo("> Do you want to replace it?") {
						continue
					}
					err := replaceProject(name, cwd, projiEnv.ConfigFolderPath, class)
					if err != nil {
						fmt.Printf("> Replacing project %s failed: %v\n", name, err)
						continue
					}
					fmt.Printf("> Project %s was successfully replaced\n", name)
				}
				continue
			}
			fmt.Printf("> Project %s was successfully created\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createProject(name, cwd, configPath string, class *models.Class) error {
	project := models.NewProject(name, filepath.Join(cwd, name), class)

	// Save it first to see if it already exists in the database
	err := projiEnv.Svc.SaveProject(project)
	if err != nil {
		return err
	}
	// Create the project
	err = project.Create(cwd, configPath)
	if err != nil {
		return err
	}
	return nil
}

func replaceProject(name, path, configPath string, class *models.Class) error {
	// Replace it
	err := projiEnv.Svc.RemoveProject(filepath.Join(path, name))
	if err != nil {
		return err
	}
	return createProject(name, path, configPath, class)
}
