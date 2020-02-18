package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kong"
	"github.com/matrix-org/gomatrix"
	"github.com/spf13/viper"
)

type cli struct {
	Config     string `short:"c" type="existingfile"`
	Homeserver string `short:"h" mapstructure:"homeserver"`
	Username   string `short:"u" mapstructure:"username"`
	Password   string `short:"p" mapstructure:"password"`

	Send sendCmd `cmd hel:"Send a message to a room"`
}

type sendCmd struct {
	Room string `arg name:"room"`
	Msg  string `arg optional name:"msg" help:"Text to send to room. Leave empty to read from stdin."`
}

func (s *sendCmd) Run(conf *cli) error {
	cl, err := gomatrix.NewClient(conf.Homeserver, "", "")
	if err != nil {
		return err
	}

	resp, err := cl.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     conf.Username,
		Password: conf.Password,
	})

	if err != nil {
		return err
	}

	cl.SetCredentials(resp.UserID, resp.AccessToken)

	msg := s.Msg

	if s.Msg == "" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		msg = string(b)
	}

	_, err = cl.SendText(s.Room, msg)

	return err
}

func main() {
	c := &cli{}
	ctx := kong.Parse(c)

	if p := c.Config; p != "" {
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

	if err := viper.Unmarshal(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := ctx.Run(&c)
	ctx.FatalIfErrorf(err)
}
