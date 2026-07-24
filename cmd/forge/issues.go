package main

import (
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage issues",
}

func init() {
	rootCmd.AddCommand(issuesCmd)
}
