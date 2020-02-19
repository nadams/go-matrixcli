package cmd

import (
	"io/ioutil"
	"os"

	"github.com/matrix-org/gomatrix"
)

type Send struct {
	Room string `arg:"" name:"room"`
	Msg  string `arg:"" optional:"" name:"msg" help:"Text to send to room. Leave empty to read from stdin."`
}

func (s *Send) Run(account Account) error {
	cl, err := gomatrix.NewClient(account.Homeserver, "", "")
	if err != nil {
		return err
	}

	resp, err := cl.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     account.Username,
		Password: account.Password,
	})

	if err != nil {
		return err
	}

	cl.SetCredentials(resp.UserID, resp.AccessToken)

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
