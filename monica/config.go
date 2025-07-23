package monica

import "github.com/spf13/viper"

type Config struct {
	Mode string `mapstructure:"mode"`
	Log  struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
}

const (
	AppModeDev  AppMode = "dev"  // development mode
	AppModeProd AppMode = "prod" // production mode
)

type AppMode string

func (m AppMode) String() string {
	return string(m)
}

func (m AppMode) IsDev() bool {
	return m == AppModeDev
}

func (m AppMode) IsProd() bool {
	return m == AppModeProd
}

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func Mode() AppMode {
	return AppMode(viper.GetString("mode"))
}
