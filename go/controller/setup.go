package controller

import (
	"fmt"
	"github.com/geobe/https-proxy/go/model"
	"github.com/spf13/viper"
)

// the relative location of project files
const Base = "src/github.com/geobe/https-proxy"

// setting up viper configuration lib
func Setup(cfgfile string) {
	if cfgfile == "" {
		cfgfile = "config"
	}
	viper.SetConfigName(cfgfile)
	viper.AddConfigPath(Base + "/config") // for config in the working directory
	viper.AddConfigPath("../../config")   // for config in the working directory
	viper.AddConfigPath(".")              // for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// read users from config file
func ReadUsers() []model.User {
	var users []model.User
	uValues := viper.Get("users").([]interface{})
	for _, uValue := range uValues {
		switch user := uValue.(type) {
		case map[string]interface{}:
			accraw := user["access"].([]interface{})
			acclist := make([]string, len(accraw))
			for index, value := range accraw {
				acclist[index] = value.(string)
			}
			nu := model.NewUser(
				user["login"].(string),
				user["password"].(string),
				acclist...)
			users = append(users, *nu)
		default:
			for k1, v1 := range uValue.(map[string]interface{}) {
				fmt.Errorf("Error in configuration filet%s: %s\n", k1, v1)
			}
		}
	}
	return users
}
