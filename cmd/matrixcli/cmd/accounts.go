package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"

	"github.com/nadams/go-matrixcli/auth"
)

type Accounts struct {
	List   ListAccounts  `cmd:"" help:"List configured accounts."`
	Login  Login         `cmd:"" help:"Log in to a matrix server."`
	Select SelectAccount `cmd:"" help:"Set the current account."`
	Remove RemoveAccount `cmd:"" help:"Remove a stored account."`
}

type ListAccounts struct {
}

func (l *ListAccounts) Run(ts *auth.TokenStore) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Homeserver", "UserID", "Current"})

	for _, a := range ts.List() {
		var current string

		if ts.CurrentName() == a.Name {
			current = "*"
		}

		t.AppendRow(table.Row{a.Name, a.Homeserver, a.UserID, current})
	}

	t.Render()

	return nil
}

type SelectAccount struct {
	Name string `arg:"" help:"Set account to be the current account."`
}

func (s *SelectAccount) Run(ts *auth.TokenStore) error {
	return ts.SetCurrent(s.Name)
}

type RemoveAccount struct {
	Name string `arg:"" help:"Account to remove."`
}

func (s *RemoveAccount) Run(ts *auth.TokenStore) error {
	return ts.Remove(s.Name)
}
