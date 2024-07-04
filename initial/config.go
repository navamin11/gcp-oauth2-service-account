package initial

import (
	"log"

	"github.com/spf13/viper"
)

var Configs *ServiceAccountConfig

func InitEnvConfigs() {
	Configs = LoadConfig()
}

type ServiceAccountConfig struct {
	ServiceAccount []struct {
		Project string `json:"project"`
		Service []struct {
			Name string `json:"name"`
			File string `json:"file"`
		} `json:"service"`
	} `json:"serviceAccount"`
}

func LoadConfig() (config *ServiceAccountConfig) {
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name

	log.Printf("%s", "Google OAuth 2.0 for Server to Server Applications")
	log.Printf("%s\n", "Copyright © 2024 By Mr.Navamin Sawasdee. All rights reserved.")

	log.Printf("%s", "-----------------------------------------------")
	log.Printf("%s", "JSON Config file")
	log.Printf("%s", "-----------------------------------------------")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("[%s] %s - %s", "X", "config.json", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("[%s] %s - %s", "X", "config.json", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("[%s] %s - %s", "X", "config.json", err)
	} else {
		log.Printf("[%s] %s", "/", "config.json")
	}
	return
}
