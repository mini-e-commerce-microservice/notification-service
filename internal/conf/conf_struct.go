package conf

type AppConfig struct {
	AppPort int    `mapstructure:"APP_PORT"`
	AppMode string `mapstructure:"APP_MODE"`
}
