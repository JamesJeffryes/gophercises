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
	"github.com/jamesjeffryes/gophercises/task/db"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, id := range args {
			db.CompleteTask(id)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		validArgs, err := db.GetTaskIDs()
		if err != nil {
			panic(err)
		}
		fmt.Println(validArgs)
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
