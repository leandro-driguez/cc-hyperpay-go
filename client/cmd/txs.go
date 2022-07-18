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

// txsCmd represents the txs command
var txsCmd = &cobra.Command{
	Use:   "txs",
	Short: "Returns all transactions involving given account",
	Long: `Returns all transactions involving given account
			receives an account id and gives the transaction history of the given account.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		contract, err := client.NewHyperPayContract()
		if err != nil {
			log.Fatalf("Failed to create contract client: %v", err)
		}
		log.Println("--> Evaluate Transaction: GetAllTxs, function gets transaction history of the given account")
		records, err := contract.Txs(id)
		if err != nil {
			log.Fatalf("Failed to evaluate transaction: %v", err)
		}
		log.Println("History of " + id + ": ")
		for i := 0; i < len(records); i++ {
			log.Println(records[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(txsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
