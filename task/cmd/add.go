/*
Copyright © 2021 James Jeffryes

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jamesjeffryes/gophercises/task/db"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		name := strings.Join(args, " ")
		dateStr, _ := cmd.Flags().GetString("duedate")
		date, err := time.Parse("01/02/2006", dateStr)
		if err != nil {
			log.Fatalln(err)
		}
		db.AddTask(name, date)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// TODO: Add Due Date
	addCmd.Flags().StringP("duedate", "d", "", "Add a date the task is due")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().Stridong("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
