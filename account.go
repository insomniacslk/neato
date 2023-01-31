package neato

import (
	"fmt"
)

func NewAccount(session *PasswordSession) *Account {
	return &Account{
		session: session,
	}
}

type Account struct {
	session *PasswordSession
	robots  []*Robot
}

func (a *Account) Robots() ([]*Robot, error) {
	if a.robots != nil {
		return a.robots, nil
	}
	var resp []*Robot
	if err := a.session.get("users/me/robots", &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch robots: %w", err)
	}
	for _, r := range resp {
		r.session = a.session
	}
	return resp, nil
}

func (a *Account) Maps() ([]*Map, error) {
	robots, err := a.Robots()
	if err != nil {
		return nil, fmt.Errorf("failed to get robots")
	}
	allMaps := make([]*Map, 0)
	for _, robot := range robots {
		maps, err := robot.Maps()
		if err != nil {
			return nil, fmt.Errorf("failed to get maps for robot '%s': %w", robot.Serial, err)
		}
		allMaps = append(allMaps, maps...)
	}
	return allMaps, nil
}
