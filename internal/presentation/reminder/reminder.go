package reminder

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/thnthien/great-deku/container"
	"github.com/thnthien/great-deku/l"
	cerrors "github.com/thnthien/great-plateau/errors"

	"telegram-reminder/config"
	reminderservice "telegram-reminder/internal/service/reminder"
)

type Config struct {
	Name     string   `json:"name"`
	Spec     string   `json:"spec,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

type Reminder struct {
	ll l.Logger `container:"name"`

	reminderService reminderservice.IService `container:"name"`

	chanID int64
}

func New(cfg *config.Config) *Reminder {
	r := &Reminder{
		chanID: cfg.RemindingChannelID,
	}
	container.Fill(r)

	return r
}

func (r *Reminder) Start(cronDelay, cronSkip *cron.Cron) {
	ctx := context.Background()
	if err := r.register(ctx, cronSkip, "0 30 11 * * *", "reminding Fe", r.reminderService.RemindFe); err != nil {
		r.ll.Fatal("cannot register RemindFe", l.Error(err))
	}
}

func (r *Reminder) register(ctx context.Context, cronManager *cron.Cron, spec string, name string, f func(ctx context.Context, chanID int64) cerrors.IError) error {
	_, err := cronManager.AddFunc(spec, func() {
		err := f(ctx, r.chanID)
		if err != nil {
			r.ll.Error("failed to remind", l.String("name", name), l.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("AddFunc: %w", err)
	}
	return nil
}
