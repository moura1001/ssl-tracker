package db_domain

import (
	"context"
	"fmt"
	"strconv"

	db_service "github.com/moura1001/ssl-tracker/src/pkg/db/service"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
)

type DomainInMemoryStore struct {
	domains  []data.DomainTracking
	domainId int64
}

func NewDomainInMemoryStore() *DomainInMemoryStore {
	return &DomainInMemoryStore{
		domains:  []data.DomainTracking{},
		domainId: 0,
	}
}

func (dbs DomainInMemoryStore) CountUserDomainTrackings(userId string) (int, error) {
	count := 0
	for _, domain := range dbs.domains {
		if domain.UserId == userId {
			count++
		}
	}
	return count, nil
}

func (dbs DomainInMemoryStore) GetDomainTrackings(filter util.Map, limit int, page int) ([]data.DomainTracking, error) {
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

func (dbs DomainInMemoryStore) GetDomainTracking(query util.Map) (*data.DomainTracking, error) {
	for _, domain := range dbs.domains {
		if isQueryMatch(domain, query) {
			return &domain, nil
		}
	}
	return nil, fmt.Errorf("no domain was found for the query %v", query)
}

func (dbs DomainInMemoryStore) GetAllTrackingsWithAccount() ([]data.TrackingAndAccount, error) {
	if len(dbs.domains) <= 0 {
		return nil, nil
	}

	trackings := make([]data.TrackingAndAccount, len(dbs.domains))
	for i := range dbs.domains {
		trackings[i] = data.TrackingAndAccount{
			NotifyUpfront:  7,
			DomainTracking: dbs.domains[i],
		}
	}

	return trackings, nil
}

func (dbs *DomainInMemoryStore) CreateDomainTrackings(trackings []data.DomainTracking) error {
	for _, domain := range trackings {
		domain.Id = dbs.domainId

		dbs.domains = append(dbs.domains, domain)

		dbs.domainId = dbs.domainId + 1
	}

	return nil
}

func (dbs *DomainInMemoryStore) UpdateAllTrackings(trackings []data.DomainTracking) error {
	for _, domain1 := range trackings {
		for i, domain2 := range dbs.domains {
			if isEquals(domain1, domain2) {
				domain1.Id = domain2.Id
				dbs.domains[i] = domain1
				break
			}
		}
	}

	return nil
}

func (dbs *DomainInMemoryStore) DeleteDomainTracking(query util.Map) error {
	for i, domain := range dbs.domains {
		if isQueryMatch(domain, query) {
			dbs.domains = append(dbs.domains[:i], dbs.domains[i+1:]...)
			return nil
		}
	}

	return nil
}

func isQueryMatch(domain data.DomainTracking, query util.Map) bool {
	isEqualsCount := 0
	queryEquals := util.Map{}
	for k, v := range query {
		if v != "" {
			switch k {
			case "id":
				isEquals := strconv.Itoa(int(domain.Id)) == v.(string)
				queryEquals[k] = isEquals
				if isEquals {
					isEqualsCount++
				}
			case "user_id":
				isEquals := domain.UserId == v.(string)
				queryEquals[k] = isEquals
				if isEquals {
					isEqualsCount++
				}
			default:
				break
			}
		}
	}

	return len(queryEquals) == isEqualsCount
}

func isEquals(domain1 data.DomainTracking, domain2 data.DomainTracking) bool {
	return domain1.UserId == domain2.UserId && domain1.DomainName == domain2.DomainName
}
