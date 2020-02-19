package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/spf13/viper"

	"github.com/nadams/go-matrixcli/cmd"
)

func main() {
	c := &cmd.CLI{}
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

	config := &cmd.Config{}

	if err := viper.Unmarshal(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(config.Accounts) == 0 {
		fmt.Println("no accounts configured!")
		os.Exit(1)
	}

	account := config.Accounts[0]

	if len(c.Account) > 0 {
		for _, a := range config.Accounts {
			if a.Name == c.Account {
				account = a
				break
			}
		}
	}

	err := ctx.Run(account)
	ctx.FatalIfErrorf(err)
}
