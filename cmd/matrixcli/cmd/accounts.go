package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"

	"github.com/nadams/go-matrixcli/auth"
)

type Accounts struct {
	List  ListAccounts `cmd:"" help:"List configured accounts."`
	Login Login        `cmd:"" help:"Log in to a matrix server."`
}

type ListAccounts struct {
}

func (l *ListAccounts) Run(ts *auth.TokenStore) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Homeserver", "UserID"})

	for _, a := range ts.List() {
		t.AppendRow(table.Row{a.Name, a.Homeserver, a.UserID})
	}

	t.Render()

	return nil
}
