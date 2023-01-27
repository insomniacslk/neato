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
	maps    []*Map
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
	if a.maps != nil {
		return a.maps, nil
	}
	robots, err := a.Robots()
	if err != nil {
		return nil, fmt.Errorf("failed to get robots")
	}
	type mapsResponse struct {
		// TODO figure out the stats field
		Stats interface{}
		Maps  []*Map
	}
	maps := make([]*Map, 0)
	var resp mapsResponse
	for _, robot := range robots {
		if err := a.session.get("users/me/robots/"+robot.Serial+"/maps", &resp); err != nil {
			return nil, fmt.Errorf("failed to get robots: %w", err)
		}
		maps = append(maps, resp.Maps...)
	}
	return maps, nil
}
