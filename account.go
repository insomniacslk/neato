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
		return nil, fmt.Errorf("failed to get robots: %w", err)
	}
	for _, r := range resp {
		r.session = a.session
	}
	return resp, nil
}
