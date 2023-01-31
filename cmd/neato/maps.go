package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagMapsShowAll bool
)

var mapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Show the most recent map of every robot",
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
		for _, r := range robots {
			maps, err := r.Maps()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get map for robot '%s' (serial: '%s'): %v\n", r.Name, r.Serial, err)
				continue
			}
			if flagJSON {
				j, err := json.Marshal(maps)
				if err != nil {
					log.Fatalf("Failed to marshal to JSON: %v", err)
				}
				fmt.Println(string(j))
			} else {
				fmt.Printf("Robot '%s' (serial: '%s')\n", r.Name, r.Serial)
				for idx, m := range maps {
					fmt.Printf("  %d) %s\n", idx+1, m)
					if !flagMapsShowAll {
						break
					}
				}
			}
		}
	},
}

func initMapsCmd() {
	mapsCmd.Flags().BoolVarP(&flagMapsShowAll, "--show-all", "a", false, "Show all the maps for each robot instead of the most recent one")
}
