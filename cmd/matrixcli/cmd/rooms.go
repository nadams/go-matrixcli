package cmd

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/matrix-org/gomatrix"

	"github.com/nadams/go-matrixcli/auth"
)

type Rooms struct {
	List ListRooms `cmd:"" help:"List joined rooms."`
}

type ListRooms struct {
	Account string `optional:"" short:"a" help:"Which account to send from."`
}

func (l *ListRooms) Run(ts *auth.TokenStore) error {
	var aa auth.AccountAuth
	var err error

	if l.Account == "" {
		aa, err = ts.Current()
	} else {
		aa, err = ts.Find(l.Account)
	}

	if err != nil {
		return err
	}
	cl, err := gomatrix.NewClient(aa.Homeserver, aa.UserID, aa.Token)
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
