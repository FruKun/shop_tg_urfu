package config

import "os"

type Config struct {
	BotToken      string
	DBPath        string
	AdminLogin    string
	AdminPassword string
	WebUrl        string
	WebPort       string
	UploadDir     string
}

func New() *Config {
	return &Config{
		BotToken:      os.Getenv("TG_BOT_TOKEN"),
		DBPath:        os.Getenv("DBPATH"),
		AdminLogin:    os.Getenv("ADMINLOGIN"),
		AdminPassword: os.Getenv("ADMINPASSWORD"),
		WebUrl:        os.Getenv("WEBURL"),
		WebPort:       os.Getenv("WEBPORT"),
		UploadDir:     os.Getenv("UPLOADDIR"),
	}
}
