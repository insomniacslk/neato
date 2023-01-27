package neato

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

// TODO implement PasswordLessSession and OAuthSession

func NewPasswordSession(endpoint string, header *url.Values) *PasswordSession {
	return &PasswordSession{
		endpoint: endpoint,
		header:   header,
	}
}

type PasswordSession struct {
	endpoint string
	header   *url.Values
}

func (s *PasswordSession) Login(email, password string) error {
	uri := "sessions"
	randBytes := make([]byte, 64)
	if _, err := rand.Read(randBytes); err != nil {
		return fmt.Errorf("failed to get random bytes")
	}

	data := map[string]interface{}{
		"email":    email,
		"password": password,
		"platform": "ios",
		"token":    hex.EncodeToString(randBytes),
	}
	type loginResponse struct {
		AccessToken string `json:"access_token"`
		CurrentTime string `json:"current_time"`
	}
	var resp loginResponse
	if err := s.post(uri, data, &resp); err != nil {
		return fmt.Errorf("http post failed: %w", err)
	}
	if s.header == nil {
		s.header = &url.Values{}
	}
	s.header.Set("Authorization", fmt.Sprintf("Token token=%s", resp.AccessToken))
	return nil
}

func (s *PasswordSession) SaveConfig() error {
	viper.Set("session.endpoint", s.endpoint)
	viper.Set("session.header", s.header)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", viper.ConfigFileUsed(), err)
	}
	return nil
}

func (s *PasswordSession) post(path string, dataMap map[string]interface{}, response interface{}) error {
	return httpPost(s.endpoint+"/"+path, s.header, dataMap, false, response)
}

func (s *PasswordSession) get(path string, response interface{}) error {
	return httpGet(s.endpoint+"/"+path, s.header, false, response)
}
