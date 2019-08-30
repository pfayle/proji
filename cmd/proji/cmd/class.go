package cmd

import (
	"github.com/spf13/cobra"
)

var classCmd = &cobra.Command{
	Use:   "class",
	Short: "Work on proji classes",
}

func init() {
	rootCmd.AddCommand(classCmd)
}