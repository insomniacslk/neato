package main

import (
	"fmt"
	"log"
	"path"

	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const progname = "neato"

var (
	defaultConfigFile = path.Join(configdir.LocalConfig(progname), "config.yml")

	flagConfigFile string
	flagToken      string
	flagDebug      bool
)

var rootCmd = &cobra.Command{
	Use:              progname,
	Short:            "neato is a CLI for BotVac's Neato cloud API",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", defaultConfigFile, "Authentication token")
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "Authentication token")
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "D", false, "Show debug output")

	// flag-name to config-directive mapping
	flagMapping := map[string]string{
		"token": "token",
	}
	for flagName, configDirective := range flagMapping {
		if err := viper.BindPFlag(configDirective, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
			log.Fatalf("Failed to bind flag --%s to config directive %s: %v", flagName, configDirective, err)
		}
	}

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(robotsCmd)
	rootCmd.AddCommand(mapsCmd)
	rootCmd.AddCommand(startCmd)
	initLoginCmd()
	initRobotsCmd()
	initMapsCmd()
	initStartCmd()
}

func initConfig() {
	viper.SetConfigFile(flagConfigFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file '%s'\n", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
