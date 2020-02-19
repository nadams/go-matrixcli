package cmd

type CLI struct {
	Account    string `optional:"" help:"Which account to use from the config file. If omitted the first one will be used."`
	ConfigFile string `optional:"" type:"existingfile" help:"Specify a config file instead of looking in default locations."`

	Send Send `cmd:"" help:"Send a message."`
}

type Config struct {
	Accounts []Account `mapstructure:"accounts"`
}

type Account struct {
	Name       string `mapstructure:"name"`
	Homeserver string `mapstructure:"homeserver"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
}
