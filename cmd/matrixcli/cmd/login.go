package cmd

import (
	"github.com/nadams/go-matrixcli/auth"
)

type Login struct {
	Name       string `help:"" short:"n" group:"login"`
	Homeserver string `help:"" short:"h" group:"login"`
	Username   string `help:"" short:"u" group:"login"`
	Password   string `help:"" short:"p" group:"login"`
}

func (l *Login) Run(ts *auth.TokenStore) error {
	_, err := ts.Login(l.Name, l.Homeserver, l.Username, l.Password)

	return err
}
