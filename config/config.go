package config

type Config struct {
	Accounts []Account `mapstructure:"accounts"`
	CacheDir string    `mapstructure:"cache_dir"`
}

type Account struct {
	Name       string `mapstructure:"name"`
	Homeserver string `mapstructure:"homeserver"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
}
