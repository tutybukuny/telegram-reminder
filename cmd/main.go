package main

import (
	"github.com/spf13/cobra"
	"github.com/thnthien/great-deku/container"
	"github.com/thnthien/great-deku/l"
	"github.com/thnthien/great-deku/l/sentry"

	"telegram-reminder/cmd/reminder"
	"telegram-reminder/config"
)

func main() {
	ll := l.New()
	cfg := config.Load(ll)

	if cfg.SentryConfig.Enabled {
		ll = l.NewWithSentry(&sentry.Configuration{
			DSN: cfg.SentryConfig.DNS,
			Trace: struct{ Disabled bool }{
				Disabled: !cfg.SentryConfig.Trace,
			},
		})
	}

	container.NamedSingleton("ll", func() l.Logger {
		return ll
	})

	cmd := newRootCmd(cfg)
	if err := cmd.Execute(); err != nil {
		ll.Error("command returns error", l.Error(err))
	}
}

func newRootCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "tg",
	}

	cmd.AddCommand(reminder.NewReminderCmd(cfg))

	return cmd
}
