package main

import (
	"log"

	"github.com/insomniacslk/neato"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagLoginCode        string
	flagLoginInteractive bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the Neato cloud",
	Run: func(cmd *cobra.Command, args []string) {
		email := viper.GetString("email")
		password := viper.GetString("password")
		if email == "" || password == "" {
			log.Fatalf("Email and password must be set")
		}
		s := neato.NewPasswordSession("https://beehive.neatocloud.com", nil)
		if err := s.Login(email, password); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
		if err := s.SaveConfig(); err != nil {
			log.Fatalf("Failed to save to config file: %v", err)
		}
		log.Printf("Saved session to config file '%s'", viper.ConfigFileUsed())
	},
}

func initLoginCmd() {
	var (
		flagLoginEmail    string
		flagLoginPassword string
	)

	loginCmd.Flags().StringVarP(&flagLoginEmail, "email", "e", "", "Email address of the Neato account")
	loginCmd.Flags().StringVarP(&flagLoginCode, "code", "C", "", "Verification code that is sent to your e-mail")
	loginCmd.Flags().StringVarP(&flagLoginPassword, "password", "p", "", "Neato account password")
	loginCmd.Flags().BoolVarP(&flagLoginInteractive, "interactive", "i", false, "Interactive login")

	flagMapping := map[string]string{
		"email":    "email",
		"password": "password",
	}
	for flagName, configDirective := range flagMapping {
		if err := viper.BindPFlag(configDirective, loginCmd.Flags().Lookup(flagName)); err != nil {
			log.Fatalf("Failed to bind flag --%s to config directive %s: %v", flagName, configDirective, err)
		}
	}
}
