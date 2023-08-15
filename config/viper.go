package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	// Run environment
	GO_ENV string `mapstructure:"GO_ENV"`

	// Mongodb config
	MONGODB_URI  string `mapstructure:"MONGODB_URI"`
	MONGODB_NAME string `mapstructure:"MONGODB_NAME"`
	PORT         string `mapstructure:"PORT"`

	// Quatalog course data urls
	COURSE_DATA_URL        string `mapstructure:"COURSE_DATA_URL"`
	COURSE_PREREQ_DATA_URL string `mapstructure:"COURSE_PREREQ_DATA_URL"`

	// Auth0 config
	AUTH0_DOMAIN   string `mapstructure:"AUTH0_DOMAIN"`
	AUTH0_AUDIENCE string `mapstructure:"AUTH0_AUDIENCE"`

	// Google oAuth config
	GOOGLE_CLIENT_ID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	OAUTH_STATE_STRING   string `mapstructure:"OAUTH_STATE_STRING"`
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
			AUTH0_DOMAIN:           os.Getenv("AUTH0_DOMAIN"),
			AUTH0_AUDIENCE:         os.Getenv("AUTH0_AUDIENCE"),
			GOOGLE_CLIENT_ID:       os.Getenv("GOOGLE_CLIENT_ID"),
			GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
			OAUTH_STATE_STRING:     os.Getenv("OAUTH_STATE_STRING"),
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

	if config.AUTH0_DOMAIN == "" {
		err = errors.New("AUTH0_DOMAIN is required")
		return
	}

	if config.AUTH0_AUDIENCE == "" {
		err = errors.New("AUTH0_AUDIENCE is required")
		return
	}

	if config.GOOGLE_CLIENT_ID == "" {
		err = errors.New("GOOGLE_CLIENT_ID is required")
		return
	}

	if config.GOOGLE_CLIENT_SECRET == "" {
		err = errors.New("GOOGLE_CLIENT_SECRET is required")
		return
	}

	if config.GO_ENV == "" {
		config.GO_ENV = "dev"
	}
	return
}
