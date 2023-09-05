package db_domain

import (
	"fmt"
	"strconv"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
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

func (dbs DomainInMemoryStore) CountDomainTrackings(filter util.Map) (int, error) {
	count := 0
	for _, domain := range dbs.domains {
		if isQueryMatch(domain, filter) {
			count++
		}
	}
	return count, nil
}

func (dbs DomainInMemoryStore) GetDomainTrackings(filter util.Map, limit int, page int) (int, []data.DomainTracking, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	offset := (limit - 1) * page

	domains := []data.DomainTracking{}

	for i, idx := offset, 0; i < len(dbs.domains) && idx < limit; {
		domain := dbs.domains[i]
		if isQueryMatch(domain, filter) {
			domains = append(domains, dbs.domains[i])
			idx++
		}

		i++
	}

	isStatusAll := filter["status"] == "all"
	if !isStatusAll && len(domains) >= limit {
		count, _ := dbs.CountDomainTrackings(filter)
		return count, domains, nil
	}

	return len(domains), domains, nil
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

func (dbs *DomainInMemoryStore) UpdateAllTrackings(trackings []*data.DomainTracking) error {
	for _, domain1 := range trackings {
		for i, domain2 := range dbs.domains {
			if isEquals(*domain1, domain2) {
				domain1.Id = domain2.Id
				dbs.domains[i] = *domain1
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
			case "domain_name":
				isEquals := domain.DomainName == v.(string)
				queryEquals[k] = isEquals
				if isEquals {
					isEqualsCount++
				}
			case "status":
				isEquals := domain.Status == v.(string)
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
