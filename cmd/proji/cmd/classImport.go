package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var remoteRepos, directories, configs, exclude []string

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 && len(remoteRepos) < 1 {
			return fmt.Errorf("no flag was passed. You have to pass the '--config', 'remote-repo' or '--directory' flag at least once")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Import configs
		// Concat the two arrays so that '... import --config *.toml' is a valid command.
		// Without appending the args, proji would only use the first toml-file and not all of
		// them as intended with the '*'.
		// TODO: This section should be optimized and cleaned up.
		for _, config := range append(configs, args...) {
			if helper.IsInSlice(exclude, config) {
				continue
			}

			err := importClassFromConfig(config, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", config, err)
				continue
			}
			fmt.Printf("> '%s' was successfully imported\n", config)
		}

		// Import directories
		for _, directory := range directories {
			confName, err := importClassFromDirectory(directory, exclude)
			if err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", directory, err)
				continue
			}
			fmt.Printf("> Directory '%s' was successfully exported to '%s'\n", directory, confName)
		}

		// Import repos
		for _, repo := range remoteRepos {
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Importing repo "
			s.Start()
			confName, err := importClassFromRemoteRepo(repo, exclude)
			s.Stop()
			if err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", repo, err)
				continue
			}
			fmt.Printf("> Repository '%s' was successfully exported to '%s'\n", repo, confName)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)

	classImportCmd.Flags().StringSliceVar(&remoteRepos, "remote-repo", []string{}, "import/imitate an existing remote repo")
	_ = classImportCmd.MarkFlagDirname("remote-repo")

	classImportCmd.Flags().StringSliceVar(&directories, "directory", []string{}, "import/imitate an existing directory")
	_ = classImportCmd.MarkFlagDirname("directory")

	classImportCmd.Flags().StringSliceVar(&configs, "config", []string{}, "import a class from a config file")
	_ = classImportCmd.MarkFlagFilename("config")

	classImportCmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "files/folders to exclude from import")
	_ = classImportCmd.MarkFlagFilename("exclude")
}

func importClassFromConfig(config string, svc storage.Service) error {
	// Import class data
	class := item.NewClass("", "", false)
	err := class.ImportFromConfig(config)
	if err != nil {
		return err
	}
	return svc.SaveClass(class)
}

func importClassFromDirectory(directory string, excludeDir []string) (string, error) {
	// Import class data
	class := item.NewClass("", "", false)
	err := class.ImportFromDirectory(directory, excludeDir)
	if err != nil {
		return "", err
	}
	return class.Export(".")
}

func importClassFromRemoteRepo(URL string, excludeDir []string) (string, error) {
	// Import class data
	class := item.NewClass("", "", false)
	err := class.ImportFromURL(URL, excludeDir)
	if err != nil {
		return "", err
	}
	return class.Export(".")
}
