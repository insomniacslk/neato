package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var mapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Show the maps of every robot",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := getAccount()
		if err != nil {
			log.Fatalf("Account lookup failed: %v", err)
		}
		maps, err := acc.Maps()
		if err != nil {
			log.Fatalf("Cannot get maps: %v", err)
		}
		if len(maps) == 0 {
			fmt.Println("No maps found")
			return
		}
		for idx, r := range maps {
			fmt.Printf("%d) %s\n", idx+1, r)
		}
	},
}

func initMapsCmd() {
}
