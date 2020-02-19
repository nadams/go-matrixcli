package cmd

import (
	"io/ioutil"
	"os"

	"github.com/matrix-org/gomatrix"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/config"
)

type Send struct {
	Room string `arg:"" name:"room"`
	Msg  string `arg:"" optional:"" name:"msg" help:"Text to send to room. Leave empty to read from stdin."`
}

func (s *Send) Run(ts *auth.TokenStore, account config.Account) error {
	aa, err := ts.Token(account.Username)
	if err != nil {
		return err
	}

	cl, err := gomatrix.NewClient(account.Homeserver, aa.UserID, aa.Token)
	if err != nil {
		return err
	}

	msg := s.Msg

	if s.Msg == "" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		msg = string(b)
	}

	_, err = cl.SendText(s.Room, msg)

	return err
}
