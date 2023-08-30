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

func (dbs DomainBunStore) GetDomainTrackings(filter util.Map, limit int, page int) ([]data.DomainTracking, error) {
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
	err := builder.Scan(context.Background())
	return trackings, err
}

// TODO: implementation
func (dbs DomainBunStore) GetDomainTracking(query util.Map) (*data.DomainTracking, error) {
	return nil, nil
}

// TODO: implementation
func (dbs DomainBunStore) GetAllTrackingsWithAccount() ([]data.TrackingAndAccount, error) {
	return []data.TrackingAndAccount{
		{DomainTracking: data.DomainTracking{DomainName: "google.com"}},
		{DomainTracking: data.DomainTracking{DomainName: "facebook.com"}},
		{DomainTracking: data.DomainTracking{DomainName: "youtube.com"}},
		{DomainTracking: data.DomainTracking{DomainName: "twitter.com"}},
		{DomainTracking: data.DomainTracking{DomainName: "amazon.com"}},
	}, nil
}

// TODO: implementation
func (dbs DomainBunStore) CreateDomainTrackings(trackings []data.DomainTracking) error {
	return nil
}

// TODO: implementation
func (dbs DomainBunStore) UpdateAllTrackings(trackings []data.DomainTracking) error {
	return nil
}

// TODO: implementation
func (dbs DomainBunStore) DeleteDomainTracking(query util.Map) error {
	return nil
}
