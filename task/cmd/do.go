package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse argument: ", arg)
			} else {
				if id <= 0 {
					fmt.Println("Invalid Task Id: ", id)
					continue
				}
				ids = append(ids, id)
			}
		}
		tasks, err := store.GetTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			return
		}
		for _, id := range ids {
			if id > len(tasks) {
				fmt.Println("Invalid Task Id: ", id)
				continue
			}
			task := tasks[id-1]
			_, err = store.CompleteTask(task.ID)
			if err != nil {
				fmt.Printf("Failed to mark '%d' as done. Error : %s\n", id, err)
				return
			}
			fmt.Printf("Marked '%d' as done\n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
