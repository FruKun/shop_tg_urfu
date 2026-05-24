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
		DBPath:        "products.db",
		AdminLogin:    "admin",
		AdminPassword: "1234",
		WebUrl:        "http://188.226.62.78:8080",
		WebPort:       ":8080",
		UploadDir:     "./uploads",
	}
}
