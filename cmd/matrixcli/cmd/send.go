package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/matrix-org/gomatrix"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/matrixext"
)

type Send struct {
	Account string `optional:"" short:"a" help:"Which account to send from."`
	Title   string `help:"Use rich formatting. If used, msg will be wrapped in a <blockquote> tag."`
	Room    string `arg:"" name:"room" help:"Can be a room ID or a room alias."`
	Msg     string `arg:"" optional:"" name:"msg" help:"Text to send to room. Leave empty to read from stdin."`
}

func (s *Send) Run(ts *auth.TokenStore) error {
	var aa auth.AccountAuth
	var err error

	if s.Account == "" {
		aa, err = ts.Current()
	} else {
		aa, err = ts.Find(s.Account)
	}

	if err != nil {
		return err
	}

	cl, err := gomatrix.NewClient(aa.Homeserver, aa.UserID, aa.Token)
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

	room := s.Room

	if !strings.Contains(room, ":") {
		domain, err := aa.Domain()
		if err != nil {
			return err
		}

		room = fmt.Sprintf("%s:%s", s.Room, domain)
	}

	if strings.HasPrefix(room, "#") {
		r, err := matrixext.GetRoomByAlias(cl, room)
		if err != nil {
			return fmt.Errorf("could not resolve room alias: %w", err)
		}

		room = r.RoomID
	}

	if s.Title != "" {
		_, err = cl.SendMessageEvent(room, "m.room.message", gomatrix.HTMLMessage{
			Body:    msg,
			MsgType: "m.text",
			Format:  "org.matrix.custom.html",
			FormattedBody: fmt.Sprintf(`
				<html>
					<body>
						<h1>%s</h1>
						<pre>%s</pre>
					</body>
				</html>`,
				s.Title,
				msg,
			),
		})
	} else {
		_, err = cl.SendText(room, msg)
	}

	return err
}
