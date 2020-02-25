package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nadams/go-matrixcli/auth"
)

type Login struct {
	Homeserver string `arg:"" help:"The homeserver to log into. With protocol."`

	Name     string `help:"The name to identify this account by." short:"n" group:"login"`
	Username string `help:"The account username." short:"u" group:"login"`
	Password string `help:"The account password." short:"p" group:"login"`
}

func (l *Login) Run(ts *auth.TokenStore) error {
	reader := bufio.NewReader(os.Stdin)

	if l.Name == "" {
		if err := l.readString(reader, "Name", &l.Name); err != nil {
			return err
		}
	}

	if l.Username == "" {
		if err := l.readString(reader, "Username", &l.Username); err != nil {
			return err
		}
	}

	if l.Password == "" {
		if err := l.readString(reader, "Password", &l.Password); err != nil {
			return err
		}
	}

	_, err := ts.Login(l.Name, l.Homeserver, l.Username, l.Password)

	return err
}

func (*Login) readString(r *bufio.Reader, prompt string, output *string) error {
	fmt.Printf("%s: ", prompt)

	s, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	*output = strings.TrimSpace(s)

	return nil
}
