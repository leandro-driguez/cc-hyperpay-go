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

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfers the given amount from the given source account to the given destination account",
	Long: `"Transfers the given amount from the given source account to the given destination account
			receives source, destination and amount, and executes the transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		dest := args[1]
		var amount float32
		_, err := fmt.Sscan(args[2], &amount)
		if err != nil {
			panic(err)
		}
		contract, err := client.NewHyperPayContract()
		if err != nil {
			log.Fatalf("Failed to create contract client: %v", err)
		}
		log.Println("--> Submit Transaction: Transfer, function transfers funds from one account to another")
		if err := contract.Transfer(source, dest, amount); err != nil {
			log.Fatalf("Failed to submit transaction: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transferCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
