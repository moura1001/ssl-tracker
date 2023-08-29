package data

import (
	"context"
	"time"

	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
)

const (
	domainTrackingTable = "domain_trackings"
	defaultLimit        = 20
)

type DomainTrackingInfo struct {
	Issuer        string
	SignatureAlgo string
	PublicKeyAlgo string
	EncodedPEM    string
	PublicKey     string
	Signature     string
	DNSNames      string
	KeyUsage      string
	ExtKeyUsages  []string `bun:",array"`
	Expires       time.Time
	Status        string
	LastPollAt    time.Time
	Latency       int
	Error         string
}

type DomainTracking struct {
	Id         int64 `bun:"id,pk,autoincrement"`
	UserId     string
	DomainName string

	DomainTrackingInfo
}

func CountUserDomainTrackings(userId string) (int, error) {
	return db.Bun.NewSelect().
		Model(&DomainTracking{}).
		Where("user_id = ?", userId).
		Count(context.Background())
}

func GetDomainTrackings(filter util.Map, limit int, page int) ([]DomainTracking, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	var trackings []DomainTracking
	builder := db.Bun.NewSelect().Model(&trackings).Limit(limit)
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
func GetDomainTracking(query util.Map) (*DomainTracking, error) {
	return nil, nil
}

// TODO: implementation
func GetAllTrackingsWithAccount() ([]TrackingAndAccount, error) {
	return []TrackingAndAccount{
		{DomainTracking: DomainTracking{DomainName: "google.com"}},
		{DomainTracking: DomainTracking{DomainName: "facebook.com"}},
		{DomainTracking: DomainTracking{DomainName: "youtube.com"}},
		{DomainTracking: DomainTracking{DomainName: "twitter.com"}},
		{DomainTracking: DomainTracking{DomainName: "amazon.com"}},
	}, nil
}

// TODO: implementation
func CreateDomainTrackings(trackings []*DomainTracking) error {
	return nil
}

// TODO: implementation
func UpdateAllTrackings(trackings []DomainTracking) error {
	return nil
}

// TODO: implementation
func DeleteDomainTracking(query util.Map) error {
	return nil
}
