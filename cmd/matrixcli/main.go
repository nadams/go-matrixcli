package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/alecthomas/kong"
	"github.com/spf13/viper"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/cmd/matrixcli/cmd"
	"github.com/nadams/go-matrixcli/config"
)

const appname = "matrixcli"

type CLI struct {
	Account    string `optional:"" help:"Which account to use from the config file. If omitted the first one will be used."`
	ConfigFile string `optional:"" type:"existingfile" help:"Specify a config file instead of looking in default locations."`
	CacheDir   string `optional:"" type:"existingdir" help:"Specify an alternate cache dir location."`

	Send      cmd.Send      `cmd:"" help:"Send a message."`
	ListRooms cmd.ListRooms `cmd:"" help:"List joined rooms"`
	Login     cmd.Login     `cmd:"" help:"Log in to a matrix server."`
}

func main() {
	c := &CLI{}
	ctx := kong.Parse(c)

	if p := c.ConfigFile; p != "" {
		viper.SetConfigFile(p)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome(), appname))
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := &config.Config{
		CacheDir: c.CacheDir,
	}

	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := initCache(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ts, err := auth.NewTokenStore(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch ctx.Command() {
	case "login":
		ctx.FatalIfErrorf(ctx.Run(ts))
	default:
		account, err := findAccount(cfg, c.Account)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ctx.Run(ts, account)
		ctx.FatalIfErrorf(err)
	}

}

func findAccount(cfg *config.Config, target string) (config.Account, error) {
	if len(cfg.Accounts) == 0 {
		return config.Account{}, errors.New("no accounts configured!")
	}

	var found bool
	account := cfg.Accounts[0]

	for _, a := range cfg.Accounts {
		if a.Name == target {
			account = a
			break
		}
	}

	if !found && target != "" {
		return config.Account{}, fmt.Errorf("could not find account %s in config\n", target)
	}

	return account, nil
}

func initCache(cfg *config.Config) error {
	if cfg.CacheDir == "" {
		cfg.CacheDir = xdg.CacheHome()

		if cfg.CacheDir == "" {
			return errors.New("could not locate cache directory")
		}

		cfg.CacheDir = filepath.Join(cfg.CacheDir, appname)

		return os.MkdirAll(cfg.CacheDir, 0755)
	}

	return nil
}
