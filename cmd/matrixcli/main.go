package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/alecthomas/kong"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/cmd/matrixcli/cmd"
)

const appname = "matrixcli"

type CLI struct {
	Account   string `optional:"" help:"Which account to use from the config file. If omitted the first one will be used."`
	ConfigDir string `optional:"" type:"existingdir" help:"Specify an alternate cache dir location."`

	Send         cmd.Send         `cmd:"" help:"Send a message."`
	ListRooms    cmd.ListRooms    `cmd:"" help:"List joined rooms"`
	Login        cmd.Login        `cmd:"" help:"Log in to a matrix server."`
	ListAccounts cmd.ListAccounts `cmd:"" help:"List configured accounts."`
}

func main() {
	c := &CLI{}
	ctx := kong.Parse(c)

	dir, err := initConfig(c.ConfigDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ts, err := auth.NewTokenStore(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := ctx.Command()

	switch {
	case strings.HasPrefix(cmd, "login"), strings.HasPrefix(cmd, "list-accounts"):
		ctx.FatalIfErrorf(ctx.Run(ts))
	default:
		account, err := findAccount(ts, c.Account)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ctx.Run(ts, account)
		ctx.FatalIfErrorf(err)
	}
}

func findAccount(ts *auth.TokenStore, target string) (auth.AccountAuth, error) {
	var a auth.AccountAuth
	var err error

	if target == "" {
		a, err = ts.First()
	} else {
		a, err = ts.Find(target)
	}

	return a, err
}

func initConfig(dir string) (string, error) {
	if dir == "" {
		dir = filepath.Join(xdg.ConfigHome(), appname)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}
