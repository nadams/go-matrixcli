package cmd

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/matrix-org/gomatrix"

	"github.com/nadams/go-matrixcli/auth"
)

type Rooms struct {
	List    ListRooms   `cmd:"" help:"List joined rooms."`
	Members ListMembers `cmd:"" help:"List members in a given room."`
}

type ListRooms struct {
	Account string `optional:"" short:"a" help:"Account to use."`
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

type ListMembers struct {
	Room    string `arg:"" help:"The room to list the members of."`
	Account string `optional:"" short:"a" help:"Account to use."`
}

func (l *ListMembers) Run(ts *auth.TokenStore) error {
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

	resp, err := cl.JoinedMembers(l.Room)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "User ID"})

	type membership struct {
		id   string
		name string
	}

	members := make([]membership, 0, len(resp.Joined))
	for id, user := range resp.Joined {
		var name string
		if n := user.DisplayName; n != nil {
			name = *n
		}

		members = append(members, membership{id: id, name: name})
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].name < members[j].name && members[i].id < members[j].id
	})

	for _, m := range members {
		t.AppendRow(table.Row{m.name, m.id})
	}

	t.Render()

	return nil
}
