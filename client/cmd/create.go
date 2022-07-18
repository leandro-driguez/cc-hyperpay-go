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
	"fmt"
	"log"

	"github.com/lllrdgz/cc-hyperpay-go/hyperpay-transfer/client"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an account with the given id, balance and bank information",
	Long: `Creates an account with the given id, balance and bank information.
			Receives id, balance and bank and create a new account with the given details.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		var balance float32
		_, err := fmt.Sscan(args[1], &balance)
		bank := args[2]
		contract, err := client.NewHyperPayContract()
		if err != nil {
			log.Fatalf("Failed to create contract client: %v", err)
		}
		log.Println("--> Submit Transaction: CreateAccount, function create a new account to the world state with given details")
		if err := contract.Create(id, balance, bank); err != nil {
			log.Fatalf("Failed to submit transaction: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
