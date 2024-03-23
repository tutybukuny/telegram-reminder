package config

import (
	"fmt"
	"strconv"

	"telegram-reminder/pkg/telegram"
)

// Config ...
type Config struct {
	Base `mapstructure:",squash"`

	MaxPoolSize  int          `json:"max_pool_size" mapstructure:"max_pool_size"`
	SentryConfig SentryConfig `json:"sentry" mapstructure:"sentry"`

	TelegramConfig telegram.Config `json:"telegram" mapstructure:"telegram"`

	RemindingChannelID int64 `json:"reminding_channel_id" mapstructure:"reminding_channel_id"`
}

// GetHTTPAddress ...
func (c *Config) GetHTTPAddress() string {
	if _, err := strconv.Atoi(c.HTTPAddress); err == nil {
		return fmt.Sprintf(":%v", c.HTTPAddress)
	}
	return c.HTTPAddress
}

// SentryConfig ...
type SentryConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled"`
	DNS     string `json:"dns" mapstructure:"dns"`
	Trace   bool   `json:"trace" mapstructure:"trace"`
}
