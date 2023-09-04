package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		start := time.Now()
		_, err := db.NewCreateTable().
			IfNotExists().
			Model((*data.Account)(nil)).
			Exec(context.Background())
		duration := time.Since(start).Milliseconds()
		CompleteTimes = append(CompleteTimes, fmt.Sprintf("%dms", duration))
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		start := time.Now()
		_, err := db.NewDropTable().
			IfExists().
			Cascade().
			Model((*data.Account)(nil)).
			Exec(context.Background())
		duration := time.Since(start).Milliseconds()
		CompleteTimes = append(CompleteTimes, fmt.Sprintf("%dms", duration))
		return err
	})
}
