package db_domain

import (
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

type DomainStore interface {
	CountUserDomainTrackings(userId string) (int, error)
	CountDomainTrackings(filter util.Map) (int, error)
	GetDomainTrackings(filter util.Map, limit int, page int) (int, []data.DomainTracking, error)
	GetDomainTracking(query util.Map) (*data.DomainTracking, error)
	GetAllTrackingsWithAccount() ([]data.TrackingAndAccount, error)
	CreateDomainTrackings(trackings []data.DomainTracking) error
	UpdateAllTrackings(trackings []*data.DomainTracking) error
	DeleteDomainTracking(query util.Map) error
}
