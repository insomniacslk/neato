package main

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cleaning",
	Run: func(cmd *cobra.Command, args []string) {
		robotIdx := 0
		if len(args) > 0 {
			n, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatalf("Invalid robot index: %v", err)
			}
			if n < 0 {
				log.Fatalf("Invalid robot index: cannot be a negative number")
			}
		}
		acc, err := getAccount()
		if err != nil {
			log.Fatalf("Account lookup failed: %v", err)
		}
		robots, err := acc.Robots()
		if err != nil {
			log.Fatalf("Cannot get robots: %v", err)
		}
		if len(robots) == 0 {
			log.Fatalf("No robots found")
		}
		if robotIdx >= len(robots) {
			log.Fatalf("Robot index is too high: got %d, must be in range 0-%d", robotIdx, len(robots)-1)
		}
		robot := robots[robotIdx]
		if err := robot.Start(nil); err != nil {
			log.Fatalf("Failed to start robot: %v", err)
		}
	},
}

func initStartCmd() {
}
