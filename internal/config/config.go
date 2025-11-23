package config

import (
	"github.com/stefanopulze/envconfig"
)

func Load() (*Config, error) {
	cfg := new(Config)
	_ = envconfig.ReadDotEnv("./.env")
	if err := envconfig.ReadEnv(cfg); err != nil {
		return nil, err
	}

	configureLogging(cfg.LogLevel, cfg.LogType)

	return cfg, nil
}

type Config struct {
	LogLevel      string `env:"LOG_LEVEL" env-default:"info"`
	LogType       string `env:"LOG_TYPE" env-default:"console"`
	Port          int    `env:"SERVER_PORT" env-default:"8080"`
	Authorization Authorization
	Providers     Providers
}

type Authorization struct {
	Username string `env:"AUTH_USERNAME"`
	Password string `env:"AUTH_PASSWORD"`
}

type Providers struct {
	Cloudflare CloudflareConfig
}

type CloudflareConfig struct {
	ApiToken string `env:"CF_API_TOKEN"`
	ZoneId   string `env:"CF_ZONE_ID"`
}
