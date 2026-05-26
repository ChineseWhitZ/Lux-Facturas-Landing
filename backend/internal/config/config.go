package config

import "os"

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	CORSOrigin  string
	LeadsFile   string
}

func Load() Config {
	return Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
		LeadsFile:   getEnv("LEADS_FILE", "data/leads.json"),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}
