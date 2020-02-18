package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/matrix-org/gomatrix"
	"github.com/spf13/viper"
)

type config struct {
	Homeserver string `mapstructure:"homeserver"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Room       string `mapstructure:"room"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/matrixcli")
	viper.AddConfigPath("$HOME/.config/matrixcli")
	viper.AddConfigPath("$HOME/.matrixcli")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var c config

	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cl, err := gomatrix.NewClient(c.Homeserver, "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp, err := cl.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     c.Username,
		Password: c.Password,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cl.SetCredentials(resp.UserID, resp.AccessToken)

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := cl.SendText(c.Room, string(b)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
