package config

import "os"

type Config struct {
	Port    string
	DBPath  string
	LogPath string
}

func Load() *Config {
	return &Config{
		Port:    getEnv("PORT", "3000"),
		DBPath:  getEnv("DB_PATH", "meaw.db"),
		LogPath: getEnv("LOG_PATH", "app.log"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
