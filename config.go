package logs

import "os"

type Config struct {
	Stage    string
	LogLevel string        `yaml:"level"`
	Debug    bool          `yaml:"debug"`
	Sentry   *SentryConfig `yaml:"sentry"`
}

type SentryConfig struct {
	Enable bool   `yaml:"enable"`
	DSN    string `yaml:"dsn"`
}

func (c *Config) SetStage() {
	c.Stage = GetEnv("STAGE", "development")
	return
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
