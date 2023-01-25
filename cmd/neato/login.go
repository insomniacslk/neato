package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
		isInteractive := flagLoginInteractive
		email := viper.GetString("email")
		code := flagLoginCode
		var err error
		if email == "" {
			if isInteractive {
				email, err = prompt("Type your Neato account email")
				if err != nil {
					log.Fatalf("Failed to read email from stdin")
				}
				email = strings.TrimSuffix(email, "\n")
			} else {
				log.Fatalf("No e-mail provided in non-interactive mode")
			}
		}
		if code == "" {
			log.Printf("No code provided, sending a login request")
			if err := loginRequestCode(email); err != nil {
				log.Fatalf("Login request failed: %v", err)
			}
			if !isInteractive {
				log.Printf("Please check your e-mail for the login code, and re-run this command with the --code flag")
				os.Exit(0)
			} else {
				code, err = prompt("Type the login code you have received via e-mail")
				if err != nil {
					log.Fatalf("Failed to read code from stdin")
				}
				code = strings.TrimSuffix(code, "\n")
			}
		}
		// get the token. At this point we have e-mail and login code
		token, err := loginRequestToken(email, code)
		if err != nil {
			log.Fatalf("Failed to get login token: %v", err)
		}
		fmt.Println(token)
	},
}

func initLoginCmd() {
	var (
		flagLoginEmail string
	)

	loginCmd.Flags().StringVarP(&flagLoginEmail, "email", "e", "", "Email address of the Neato account")
	loginCmd.Flags().StringVarP(&flagLoginCode, "code", "C", "", "Verification code that is sent to your e-mail")
	loginCmd.Flags().BoolVarP(&flagLoginInteractive, "interactive", "i", false, "Interactive login")

	flagMapping := map[string]string{
		"email": "email",
	}
	for flagName, configDirective := range flagMapping {
		if err := viper.BindPFlag(configDirective, loginCmd.Flags().Lookup(flagName)); err != nil {
			log.Fatalf("Failed to bind flag --%s to config directive %s: %v", flagName, configDirective, err)
		}
	}
}

func loginRequestCode(email string) error {
	uri := "https://mykobold.eu.auth0.com/passwordless/start"
	body := map[string]string{
		"send":       "code",
		"email":      email,
		"client_id":  "KY4YbVAvtgB7lp8vIbWQ7zLk3hssZlhR",
		"connection": "email",
	}
	// ignore response content, we assume that an HTTP 200 is good enough
	_, err := loginRequest(uri, body)
	return err
}

func loginRequestToken(email, code string) (string, error) {
	uri := "https://mykobold.eu.auth0.com/oauth/token"
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		Scope       string `json:"scope"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	body := map[string]string{
		"prompt":       "login",
		"grant_type":   "http://auth0.com/oauth/grant-type/passwordless/otp",
		"scope":        "openid email profile read:current_user",
		"locale":       "en",
		"otp":          code,
		"source":       "vorwerk_auth0",
		"platform":     "ios",
		"audience":     "https://mykobold.eu.auth0.com/userinfo",
		"username":     email,
		"client_id":    "KY4YbVAvtgB7lp8vIbWQ7zLk3hssZlhR",
		"realm":        "email",
		"country_code": "DE",
	}
	resp, err := loginRequest(uri, body)
	if err != nil {
		return "", err
	}
	var tr tokenResponse
	if err := json.Unmarshal(resp, &tr); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return tr.IDToken, nil
}

func loginRequest(uri string, dataMap map[string]string) ([]byte, error) {
	data, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP body: %w", err)
	}
	if flagDebug {
		log.Printf("login request response: %s", body)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected HTTP 200 OK, got %s", resp.Status)
	}
	return body, nil
}

func prompt(text string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text + " > ")
	return reader.ReadString('\n')
}
