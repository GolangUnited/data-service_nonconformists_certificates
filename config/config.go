package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all parameters needed
// for app launch
type Config struct {
	Port             int
	DBType           string
	ConnectionString string
}

// GetConfig reads from env vars with CERTMGR_ prefix and
// from .env file and returns Config structure
func GetConfig() (Config, error) {
	log.Println("reading config...")
	viper.SetConfigFile(".env")
	viper.SetEnvPrefix("certmgr")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	var conf Config
	if err != nil {
		return conf, err
	}
	conf.Port = viper.GetInt("port")
	conf.DBType = viper.GetString("dbtype")
	conf.ConnectionString = viper.GetString("ConnectionString")
	return conf, nil
}
