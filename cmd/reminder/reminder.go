package reminder

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/thnthien/great-deku/container"
	"github.com/thnthien/great-deku/l"
	"telegram-reminder/pkg/cronlogger"

	"telegram-reminder/config"
	"telegram-reminder/internal/presentation/reminder"
)

func NewReminderCmd(cfg *config.Config) *cobra.Command {
	var (
		cronSkip *cron.Cron
	)

	cmd := &cobra.Command{
		Use: "reminder",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var ll l.Logger
			container.NamedSingleton("ll", &ll)
			crll := cronlogger.NewCronLogger(ll)
			bootstrap(cfg)
			cronSkip = cron.New(cron.WithParser(
				cron.NewParser(
					cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
				cron.WithChain(cron.Recover(crll), cron.SkipIfStillRunning(crll)))
			cronSkip.Start()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer cronSkip.Stop()
			return nil
		},
	}

	return cmd
}

func startReminder(cfg *config.Config, cronDelay, cronSkip *cron.Cron) {
	r := reminder.New(cfg)

	r.Start(cronDelay, cronSkip)
}
