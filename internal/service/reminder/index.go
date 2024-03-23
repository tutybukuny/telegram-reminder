package reminderservice

import (
	"context"

	cerrors "github.com/thnthien/great-plateau/errors"
)

type IService interface {
	RemindFe(ctx context.Context, chanID int64) cerrors.IError
}
