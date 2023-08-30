package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/settings"
	"github.com/moura1001/ssl-tracker/src/pkg/ssl"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

var limitFilters = []int{
	5,
	10,
	25,
	50,
}

var statusFilters = []string{
	"all",
	data.StatusHealthy,
	data.StatusExpires,
	data.StatusExpired,
	data.StatusInvalid,
	data.StatusOffline,
}

func HandleDomainList(ctx *gin.Context) {
	user := getAuthenticatedUser(ctx)
	count, err := db.Store.Domain.CountUserDomainTrackings(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	if count <= 0 {
		ctx.HTML(http.StatusOK, "domains/index", util.Map{"userHasTrackings": false})
		return
	}
	filter, err := buildTrackingFilter(ctx)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	filterContext := buildFilterContext(filter)
	query := util.Map{
		"user_id": user.Id,
	}
	if filter.Status != "all" {
		query["status"] = filter.Status
	}
	domains, err := db.Store.Domain.GetAllTrackingsWithAccount()
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	data := util.Map{
		"trackings":        domains,
		"filters":          filterContext,
		"userHasTrackings": true,
		"pages":            buildPages(count, filter.Limit),
		"queryParams":      filter.encode(),
	}
	ctx.HTML(http.StatusOK, "domains/index", data)
}

func HandleDomainNew(ctx *gin.Context) {
	flashes, _ := ctx.Get("flash")
	ctx.HTML(http.StatusOK, "domains/new", util.Map{
		"flash": flashes,
	})
}

func HandleDomainDelete(ctx *gin.Context) {
	user := getAuthenticatedUser(ctx)
	query := util.Map{
		"user_id": user.Id,
		"id":      ctx.Param("id"),
	}
	if err := db.Store.Domain.DeleteDomainTracking(query); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	ctx.Redirect(http.StatusFound, "/domains")
}

func HandleDomainShow(ctx *gin.Context) {
	trackingId := ctx.Param("id")
	user := getAuthenticatedUser(ctx)
	query := util.Map{
		"user_id": user.Id,
		"id":      trackingId,
	}
	tracking, err := db.Store.Domain.GetDomainTracking(query)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	context := util.Map{
		"tracking": tracking,
	}
	ctx.HTML(http.StatusOK, "domains/show", context)
}

func HandleDomainCreate(ctx *gin.Context) {
	flashData := util.Map{}
	userDomainsInput := ctx.Request.FormValue("domains")
	userDomainsInput = strings.ReplaceAll(userDomainsInput, " ", "")

	if len(userDomainsInput) <= 0 {
		flashData["domainsError"] = "Please provide at least 1 valid domain name"
		flashWithData(ctx, flashData)
		ctx.Redirect(http.StatusFound, "/domains/new")
		return
	}
	domains := strings.Split(userDomainsInput, ",")
	if len(domains) <= 0 {
		flashData["domainsError"] = "Invalid domain list input. Make sure to use a comma separated list (domain1.com, domain2.com, etc)"
		flashData["domains"] = userDomainsInput
		flashWithData(ctx, flashData)
		ctx.Redirect(http.StatusFound, "/domains/new")
		return
	}
	for _, domain := range domains {
		if !util.IsValidDomainName(domain) {
			flashData["domainsError"] = fmt.Sprintf("%s is not a valid domain", domain)
			flashData["domains"] = userDomainsInput
			flashWithData(ctx, flashData)
			ctx.Redirect(http.StatusFound, "/domains/new")
			return
		}
	}

	user := getAuthenticatedUser(ctx)
	account, err := db.Store.Account.GetUserAccount(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}

	maxTrackings := settings.Account[account.Plan].MaxTrackings
	count, err := db.Store.Domain.CountUserDomainTrackings(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	if len(domains)+count > maxTrackings {
		flashData["maxTrackings"] = maxTrackings
		flashData["domains"] = userDomainsInput
		flashWithData(ctx, flashData)
		ctx.Redirect(http.StatusFound, "/domains/new")
		return
	}

	if err := createAllDomainTrackings(user.Id, domains); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}

	ctx.Redirect(http.StatusFound, "/domains")
}

func createAllDomainTrackings(userId string, domains []string) error {
	var (
		trackingsChan = make(chan data.DomainTracking, len(domains))
		wg            = sync.WaitGroup{}
	)

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer func() {
				cancel()
				wg.Done()
			}()

			info, err := ssl.PollDomain(ctx, domain)
			if err != nil {
				return
			}

			trackingsChan <- data.DomainTracking{
				UserId:             userId,
				DomainName:         domain,
				DomainTrackingInfo: *info,
			}
		}(domain)
	}
	wg.Wait()
	close(trackingsChan)

	return db.Store.Domain.CreateDomainTrackings(processResults(trackingsChan))
}

type trackingFilter struct {
	Limit  int
	Page   int
	Status string
	Sort   string
}

func (f *trackingFilter) encode() string {
	values := url.Values{}
	if f.Limit > 0 {
		values.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Page > 0 {
		values.Set("page", strconv.Itoa(f.Page))
	}
	values.Set("status", f.Status)
	return values.Encode()
}

func buildTrackingFilter(ctx *gin.Context) (*trackingFilter, error) {
	filter := new(trackingFilter)
	if err := ctx.ShouldBindQuery(filter); err != nil {
		return nil, err
	}
	if filter.Limit == 0 {
		filter.Limit = 25
	} else if filter.Limit < 0 {
		filter.Limit = int(math.Abs(float64(filter.Limit)))
	}
	return filter, nil
}

func buildFilterContext(filter *trackingFilter) util.Map {
	return util.Map{
		"statuses":       statusFilters,
		"limits":         limitFilters,
		"selectedStatus": filter.Status,
		"selectedLimit":  filter.Limit,
		"selectedPage":   filter.Page,
	}
}

func buildPages(results int, limit int) []int {
	pages := make([]int, results/limit)
	for i := 0; i < len(pages); i++ {
		pages[i] = i + 1
	}
	return pages
}
