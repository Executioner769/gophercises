package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopher.com/task/db"
)

var store db.Store

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI Task Manager",
}

func Execute(s db.Store) {
	store = s
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
