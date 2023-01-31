package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var robotsCmd = &cobra.Command{
	Use:   "robots",
	Short: "Show the robots under the currently logged in account",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := getAccount()
		if err != nil {
			log.Fatalf("Account lookup failed: %v", err)
		}
		robots, err := acc.Robots()
		if err != nil {
			log.Fatalf("Cannot get robots: %v", err)
		}
		if len(robots) == 0 {
			fmt.Println("No robots found")
			return
		}
		if flagJSON {
			j, err := json.Marshal(robots)
			if err != nil {
				log.Fatalf("Failed to marshal to JSON: %v", err)
			}
			fmt.Println(string(j))
		} else {
			for idx, r := range robots {
				fmt.Printf("%d) %s\n", idx+1, r)
			}
		}
	},
}

func initRobotsCmd() {
}
