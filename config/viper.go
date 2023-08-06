package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	// Mongodb config
	MONGODB_URI  string `mapstructure:"MONGODB_URI"`
	MONGODB_NAME string `mapstructure:"MONGODB_NAME"`
	PORT         string `mapstructure:"PORT"`

	// Quatalog course data urls
	COURSE_DATA_URL        string `mapstructure:"COURSE_DATA_URL"`
	COURSE_PREREQ_DATA_URL string `mapstructure:"COURSE_PREREQ_DATA_URL"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			MONGODB_URI:            os.Getenv("MONGODB_URI"),
			MONGODB_NAME:           os.Getenv("MONGODB_NAME"),
			PORT:                   os.Getenv("PORT"),
			COURSE_DATA_URL:        os.Getenv("COURSE_DATA_URL"),
			COURSE_PREREQ_DATA_URL: os.Getenv("COURSE_PREREQ_DATA_URL"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.MONGODB_URI == "" {
		err = errors.New("MONGODB_URI is required")
		return
	}

	if config.MONGODB_NAME == "" {
		err = errors.New("MONGODB_NAME is required")
		return
	}

	if config.COURSE_DATA_URL == "" {
		err = errors.New("COURSE_DATA_URL is required")
		return
	}

	if config.COURSE_DATA_URL == "" {
		err = errors.New("COURSE_PREREQ_DATA_URL is required")
		return
	}

	return
}
