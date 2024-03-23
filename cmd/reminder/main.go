package reminder

import (
	"github.com/robfig/cron/v3"
	"github.com/thnthien/great-deku/container"
	handleossignal "github.com/thnthien/great-deku/handle-os-signal"
	"github.com/thnthien/great-deku/l"
	"github.com/thnthien/great-deku/l/sentry"

	"telegram-reminder/config"
	"telegram-reminder/pkg/cronlogger"
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

	// init os signal handle
	shutdown := handleossignal.New(ll)
	shutdown.HandleDefer(func() {
		ll.Sync()
	})
	container.NamedSingleton("shutdown", func() handleossignal.IShutdownHandler {
		return shutdown
	})

	bootstrap(cfg)

	crll := cronlogger.NewCronLogger(ll)
	cronDelay := cron.New(cron.WithParser(
		cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithChain(cron.Recover(crll), cron.DelayIfStillRunning(crll)))
	cronSkip := cron.New(cron.WithParser(
		cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithChain(cron.Recover(crll), cron.SkipIfStillRunning(crll)))

	go startReminder(cfg, cronDelay, cronSkip)

	cronDelay.Start()
	defer cronDelay.Stop()
	cronSkip.Start()
	defer cronSkip.Stop()

	// handle signal
	if cfg.Environment == "D" {
		shutdown.SetTimeout(1)
	}
	shutdown.Handle()
}
