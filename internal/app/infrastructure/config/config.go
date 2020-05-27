package config

import "os"

type Config struct {
	Env            string
	TelegramApiKey string
	VkApiKey       string
	UseProxy       string
}

func New() *Config {
	cfg := &Config{
		Env:            os.Getenv("ENV"),
		TelegramApiKey: os.Getenv("TELEGRAM_API_KEY"),
		VkApiKey:       os.Getenv("VK_API_KEY"),
		UseProxy:       os.Getenv("USE_PROXY"),
	}
	return cfg
}