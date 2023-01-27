package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/insomniacslk/neato"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var robotsCmd = &cobra.Command{
	Use:   "robots",
	Short: "Show the robots under the currently logged in account",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := viper.GetString("session.endpoint")
		header := url.Values{}
		headerList := viper.Get("session.header").(map[string]interface{})
		for k, vi := range headerList {
			v := vi.([]interface{})
			for _, h := range v {
				header.Add(k, h.(string))
			}
		}
		if endpoint == "" || header.Get("authorization") == "" {
			log.Fatalf("No session.endpoint or session.header.Authorization found in configuration file, you need to log in first")
		}
		s := neato.NewPasswordSession(endpoint, &header)
		a := neato.NewAccount(s)
		robots, err := a.Robots()
		if err != nil {
			log.Fatalf("Cannot get robots: %v", err)
		}
		for idx, r := range robots {
			fmt.Printf("%d) %s\n", idx+1, r)
		}
	},
}

func initRobotsCmd() {
}
