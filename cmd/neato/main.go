package main

import (
	"fmt"
	"log"
	"os"
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
	flagJSON       bool
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
	rootCmd.PersistentFlags().BoolVarP(&flagJSON, "json", "j", false, "Print output as JSON")

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
	rootCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	initLoginCmd()
	initRobotsCmd()
	initMapsCmd()
	initStateCmd()
	initStartCmd()
	initStopCmd()
}

func initConfig() {
	viper.SetConfigFile(flagConfigFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file '%s'\n", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
