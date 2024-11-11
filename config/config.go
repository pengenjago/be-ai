package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var profileViper *viper.Viper

func SetDefault(name string, val interface{}) {

	if profileViper != nil {
		profileViper.SetDefault(name, val)
	} else {
		viper.SetDefault(name, val)
	}
}

func HasConfig(name string) bool {

	if profileViper != nil {
		return profileViper.InConfig(name)
	}

	return viper.InConfig(name)
}

func Get(name string, defVal ...string) string {

	if profileViper != nil && profileViper.InConfig(name) {
		return profileViper.GetString(name)
	} else if viper.InConfig(name) {
		return viper.GetString(name)
	} else if len(defVal) > 0 {
		return defVal[0]
	}

	return ""
}

// GetConfigEnv is First get from os.env if not found get from file config
func GetConfigEnv(key string) string {

	cfg := os.Getenv(key)
	if cfg == "" {
		cfg = Get(key)
	}

	return cfg
}

func GetInt(name string, defVal ...int) int {

	if profileViper != nil && profileViper.InConfig(name) {
		return profileViper.GetInt(name)
	} else if viper.InConfig(name) {
		return viper.GetInt(name)
	} else if len(defVal) > 0 {
		return defVal[0]
	}

	return 0
}

func GetBool(name string, defVal ...bool) bool {

	if profileViper != nil && profileViper.InConfig(name) {
		return profileViper.GetBool(name)
	} else if viper.InConfig(name) {
		return viper.GetBool(name)
	} else if len(defVal) > 0 {
		return defVal[0]
	}

	return false
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error read config : ", err)
		return
	}
}
