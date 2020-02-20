package cmd

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/matrix-org/gomatrix"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/config"
)

type ListRooms struct{}

func (l *ListRooms) Run(ts *auth.TokenStore, account config.Account) error {
	aa, err := ts.Token(account.Username)
	if err != nil {
		return err
	}

	cl, err := gomatrix.NewClient(account.Homeserver, aa.UserID, aa.Token)
	if err != nil {
		return err
	}

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
