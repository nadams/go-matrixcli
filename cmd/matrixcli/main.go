package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/alecthomas/kong"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/cmd/matrixcli/cmd"
)

const appname = "matrixcli"

type CLI struct {
	ConfigDir string `optional:"" type:"existingdir" help:"Specify an alternate cache dir location."`

	Send     cmd.Send     `cmd:"" help:"Send a message."`
	Rooms    cmd.Rooms    `cmd:"" help:"Operate on rooms."`
	Accounts cmd.Accounts `cmd:"" help:"Operate on configured accounts."`
}

func main() {
	c := &CLI{}
	ctx := kong.Parse(c)

	dir, err := initConfig(c.ConfigDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ts, err := auth.NewTokenStore(dir, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx.FatalIfErrorf(ctx.Run(ts))
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
