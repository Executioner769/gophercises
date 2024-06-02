package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	allFlag       bool
	completedFlag bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "lists all the tasks which are not marked as done",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := store.GetTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks to complete! Why not take a vacation? 🏝️")
			return
		}

		for i, task := range tasks {
			taskMsg := fmt.Sprintf("%d. %s\n", i+1, task.Message)
			if allFlag {
				if !task.Done {
					color.Magenta(taskMsg)
				} else {
					color.Cyan(taskMsg)
				}
			} else if completedFlag {
				if task.Done {
					color.Cyan(taskMsg)
				}
			} else {
				if !task.Done {
					color.Magenta(taskMsg)
				}
			}
		}
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "list all the tasks")
	lsCmd.Flags().BoolVarP(&completedFlag, "completed", "c", false, "list all the tasks that are marked as done")

	rootCmd.AddCommand(lsCmd)
}
