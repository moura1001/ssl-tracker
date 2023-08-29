package notify

import (
	"context"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
)

type Notifier interface {
	Notify(ctx context.Context, account data.TrackingAndAccount) error
	Kind() string
}
