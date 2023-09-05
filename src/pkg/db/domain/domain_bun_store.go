package db_domain

import (
	"context"

	db_service "github.com/moura1001/ssl-tracker/src/pkg/db/service"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
)

const (
	domainTrackingTable = "domain_trackings"
	defaultLimit        = 20
)

type DomainBunStore struct{}

func NewDomainBunStore() DomainBunStore {
	return DomainBunStore{}
}

func (dbs DomainBunStore) CountUserDomainTrackings(userId string) (int, error) {
	return db_service.Bun.NewSelect().
		Model(&data.DomainTracking{}).
		Where("user_id = ?", userId).
		Count(context.Background())
}

func (dbs DomainBunStore) CountDomainTrackings(filter util.Map) (int, error) {
	builder := db_service.Bun.NewSelect().Model(&data.DomainTracking{})
	for k, v := range filter {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	return builder.Count(context.Background())
}

func (dbs DomainBunStore) GetDomainTrackings(filter util.Map, limit int, page int) (int, []data.DomainTracking, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	var trackings []data.DomainTracking
	builder := db_service.Bun.NewSelect().Model(&trackings).Limit(limit)
	for k, v := range filter {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	offset := (limit - 1) * page
	builder.Offset(offset)
	err := builder.OrderExpr("domain_tracking.id ASC").
		Scan(context.Background())

	isStatusAll := filter["status"] == "all"
	if err == nil && !isStatusAll && len(trackings) >= limit {
		count, err := dbs.CountDomainTrackings(filter)
		if err == nil {
			return count, trackings, nil
		}
	}

	return len(trackings), trackings, err
}

func (dbs DomainBunStore) GetDomainTracking(query util.Map) (*data.DomainTracking, error) {
	var tracking data.DomainTracking
	builder := db_service.Bun.NewSelect().Model(&tracking)
	for k, v := range query {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	err := builder.Scan(context.Background())
	return &tracking, err
}

func (dbs DomainBunStore) GetAllTrackingsWithAccount() ([]data.TrackingAndAccount, error) {
	var trackings []data.TrackingAndAccount
	err := db_service.Bun.NewSelect().
		Table(domainTrackingTable).
		ColumnExpr("domain_trackings.*").
		ColumnExpr("a.notify_upfront").
		Join("JOIN accounts AS a ON a.user_id = domain_trackings.user_id").
		OrderExpr("domain_trackings.id ASC").
		Scan(context.Background(), &trackings)
	return trackings, err
}

func (dbs DomainBunStore) CreateDomainTrackings(trackings []data.DomainTracking) error {
	if len(trackings) > 0 {
		_, err := db_service.Bun.NewInsert().
			Model(&trackings).
			Ignore().
			Exec(context.Background())
		return err
	}
	return nil
}

func (dbs DomainBunStore) UpdateAllTrackings(trackings []*data.DomainTracking) error {
	if len(trackings) > 0 {
		values := db_service.Bun.NewValues(&trackings)
		_, err := db_service.Bun.NewUpdate().
			With("_data", values).
			Model((*data.DomainTracking)(nil)).
			TableExpr("_data").
			Set("status = _data.status").
			Set("latency = _data.latency").
			Set("last_poll_at = _data.last_poll_at").
			Where("domain_tracking.id = _data.id").
			Where("domain_tracking.user_id = _data.user_id").
			Where("domain_tracking.domain_name = _data.domain_name").
			Exec(context.Background())
		return err
	}
	return nil
}

func (dbs DomainBunStore) DeleteDomainTracking(query util.Map) error {
	builder := db_service.Bun.NewDelete().Model(&data.DomainTracking{})
	for k, v := range query {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	_, err := builder.Exec(context.Background())
	return err
}
