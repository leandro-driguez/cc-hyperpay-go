/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"log"

	"github.com/lllrdgz/cc-hyperpay-go/hyperpay-transfer/client"
	"github.com/spf13/cobra"
)

// existsCmd represents the exists command
var existsCmd = &cobra.Command{
	Use:   "exists",
	Short: "Determines whether an account with the given ID exists",
	Long: `Determines whether an account with the given ID exists.
			Receives an account id.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		contract, err := client.NewHyperPayContract()
		if err != nil {
			log.Fatalf("Failed to create contract client: %v", err)
		}
		log.Println("--> Evaluate Transaction: AccountExists, function returns true if the given account exists in the world state")
		exists, err := contract.Exists(id)
		if err != nil {
			log.Fatalf("Failed to evaluate transaction: %v", err)
		}
		if exists {
			log.Println("The account " + id + " exists")
		} else {
			log.Println("The account " + id + " doesn't exists")
		}
	},
}

func init() {
	rootCmd.AddCommand(existsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// existsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// existsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
