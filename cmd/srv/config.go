package main

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerAddr string `mapstructure:"SERVER_ADDRESS"`
	MyDB       string `mapstructure:"MYSQL_DB"`
	MyUser     string `mapstructure:"MYSQL_USER"`
	MyPassword string `mapstructure:"MYSQL_PASSWORD"`
	MyHost     string `mapstructure:"MYSQL_HOST"`
	MyPort     string `mapstructure:"MYSQL_PORT"`
	Secret     string `mapstructure:"SECRET"`
}

func LoadEnvVariables() *Config {
	var c Config

	//tell viper the path of env file
	viper.AddConfigPath("./")
	//tell viper the name of file
	viper.SetConfigName("config")
	//tell viper type of file
	viper.SetConfigType("env")

	//reads all the variables from env file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading env file", err)
	}

	//unmarshal the loaded env variables
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalln(err)
	}
	return &c
}
