package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var removeAll bool

var classRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if removeAll {
			err := removeAllClasses(projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing of all classes failed: %v\n", err)
				return err
			}
			fmt.Println("> All classes were successfully exported")
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}

		for _, name := range args {
			if err := removeClass(name, projiEnv.Svc); err != nil {
				fmt.Printf("> Removing '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("> '%s' was successfully removed\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
	classRmCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all classes")
}

func removeClass(label string, svc storage.Service) error {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}

	if classID == 1 {
		return fmt.Errorf("Class 1 can not be removed")
	}

	return svc.RemoveClass(classID)
}

func removeAllClasses(svc storage.Service) error {
	classes, err := svc.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		if class.ID == 1 {
			continue
		}
		err = svc.RemoveClass(class.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
