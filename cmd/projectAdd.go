package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add LABEL PATH",
	Short: "Add an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("missing label or path")
		}

		path, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		if !helper.DoesPathExist(path) {
			return fmt.Errorf("path '%s' does not exist", path)
		}

		label := strings.ToLower(args[0])

		err = addProject(label, path)
		if err != nil {
			return err
		}
		fmt.Printf("> Project '%s' was successfully added\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addProject(label, path string) error {
	name := filepath.Base(path)
	classID, err := projiEnv.Svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}

	class, err := projiEnv.Svc.LoadClass(classID)
	if err != nil {
		return err
	}

	proj := item.NewProject(0, name, path, class)
	return projiEnv.Svc.SaveProject(proj)
}