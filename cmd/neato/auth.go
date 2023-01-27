package main

import (
	"fmt"
	"net/url"

	"github.com/insomniacslk/neato"
	"github.com/spf13/viper"
)

func getAccount() (*neato.Account, error) {
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
		return nil, fmt.Errorf("no session.endpoint or session.header.Authorization found in configuration file, you need to log in first")
	}
	s := neato.NewPasswordSession(endpoint, &header)
	return neato.NewAccount(s), nil
}
