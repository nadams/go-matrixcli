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

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const appname = "matrixcli"

type printVersion struct{}

func (v *printVersion) Run() error {
	fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)

	return nil
}

type CLI struct {
	ConfigDir string `optional:"" type:"existingdir" help:"Specify an alternate cache dir location."`

	Send     cmd.Send     `cmd:"" help:"Send a message."`
	Rooms    cmd.Rooms    `cmd:"" help:"Operate on rooms."`
	Accounts cmd.Accounts `cmd:"" help:"Operate on configured accounts."`
	Version  printVersion `cmd:"" help:"Print version information and exit."`
}

func main() {
	c := &CLI{}
	ctx := kong.Parse(c)

	if ctx.Command() == "version" {
		ctx.FatalIfErrorf(ctx.Run())
		os.Exit(0)
	}

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
