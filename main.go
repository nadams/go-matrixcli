package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/alecthomas/kong"
	"github.com/spf13/viper"

	"github.com/nadams/go-matrixcli/auth"
	"github.com/nadams/go-matrixcli/cmd"
	"github.com/nadams/go-matrixcli/config"
)

type CLI struct {
	Account    string `optional:"" help:"Which account to use from the config file. If omitted the first one will be used."`
	ConfigFile string `optional:"" type:"existingfile" help:"Specify a config file instead of looking in default locations."`
	CacheDir   string `optional:"" type:"existingdir" help:"Specify an alternate cache dir location."`

	Send cmd.Send `cmd:"" help:"Send a message."`
}

func main() {
	c := &CLI{}
	ctx := kong.Parse(c)

	if p := c.ConfigFile; p != "" {
		viper.SetConfigFile(p)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/matrixcli")
		viper.AddConfigPath("$HOME/.config/matrixcli")
		viper.AddConfigPath("$HOME/.matrixcli")
		viper.AddConfigPath(".")
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := &config.Config{
		CacheDir: c.CacheDir,
	}

	if cfg.CacheDir == "" {
		cfg.CacheDir = xdg.CacheHome()

		if cfg.CacheDir == "" {
			fmt.Println("could not locate cache directory")
			os.Exit(1)
		}

		cfg.CacheDir = filepath.Join(cfg.CacheDir, "matrixcli")

		if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(cfg.Accounts) == 0 {
		fmt.Println("no accounts configured!")
		os.Exit(1)
	}

	account := cfg.Accounts[0]

	if len(c.Account) > 0 {
		for _, a := range cfg.Accounts {
			if a.Name == c.Account {
				account = a
				break
			}
		}
	}

	ts, err := auth.NewTokenStore(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ctx.Run(ts, account)
	ctx.FatalIfErrorf(err)
}
