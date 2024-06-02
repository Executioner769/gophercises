package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopher.com/task/db"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds a task to your tasks list",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires at least one task")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")
		task := db.NewTask(message)
		_, err := store.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			return
		}
		fmt.Printf("Added '%s' to your Tasks list.\n", message)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
