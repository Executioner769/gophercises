package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var compCmd = &cobra.Command{
	Use:        "comp",
	Short:      "lists all the completed tasks",
	Deprecated: "Instead run $ task ls -c",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := store.GetTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		for i, task := range tasks {
			if task.Done {
				fmt.Printf("%d. %s\n", i+1, task.Message)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compCmd)
}
