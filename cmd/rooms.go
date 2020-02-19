package cmd

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/matrix-org/gomatrix"
)

type ListRooms struct{}

func (l *ListRooms) Run(account Account) error {
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

	rooms, err := cl.JoinedRooms()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Room"})

	sort.Strings(rooms.JoinedRooms)

	for i, r := range rooms.JoinedRooms {
		t.AppendRow(table.Row{i + 1, r})
	}

	t.Render()

	return nil
}
