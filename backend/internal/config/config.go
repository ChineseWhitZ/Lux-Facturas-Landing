package config

import "os"

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	CORSOrigin  string
	LeadsFile   string
	AdminUser   string
	AdminPass   string
	SessionKey  string
}

func Load() Config {
	return Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
		LeadsFile:   getEnv("LEADS_FILE", "data/leads.json"),
		AdminUser:   getEnv("ADMIN_USER", "admin"),
		AdminPass:   getEnv("ADMIN_PASSWORD", "lux-admin"),
		SessionKey:  getEnv("SESSION_KEY", "dev-session-key-change-me"),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}
