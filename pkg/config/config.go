package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port             int
	DBType           string
	ConnectionString string
}

func GetConfig() (conf Config, err error) {
	log.Println("reading config...")
	viper.SetConfigFile(".env")
	viper.SetEnvPrefix("certmgr")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	if err != nil {
		return conf, err
	}
	conf.Port = viper.GetInt("port")
	conf.DBType = viper.GetString("dbtype")
	conf.ConnectionString = viper.GetString("ConnectionString")
	return conf, nil
}
