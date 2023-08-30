package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/ssl"
)

type Monitor struct {
	interval time.Duration
	lastPoll time.Time
	quitch   chan struct{}
}

func NewMonitor(interval time.Duration) *Monitor {
	return &Monitor{
		interval: interval,
		quitch:   make(chan struct{}),
	}
}

func (m *Monitor) poll() (int, error) {
	trackingsWithAccount, err := db.Store.Domain.GetAllTrackingsWithAccount()
	if err != nil {
		return 0, err
	}

	var (
		workers = make(chan struct{}, 15)
		wg      = sync.WaitGroup{}
		results = make(chan data.DomainTracking, len(trackingsWithAccount))
	)
	for _, trackingWithAccount := range trackingsWithAccount {
		wg.Add(1)
		go func(tracking data.TrackingAndAccount) {
			workers <- struct{}{}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer func() {
				<-workers
				wg.Done()
				cancel()
			}()

			domainName := tracking.DomainName
			info, err := ssl.PollDomain(ctx, domainName)
			if err != nil {
				logger.Log("error", "poll domain error", "err", err, "domain", domainName)
				return
			}

			domainTracking := &tracking.DomainTracking
			domainTracking.DomainTrackingInfo = *info

			m.maybeNotify(context.Background(), tracking)
			results <- *domainTracking
		}(trackingWithAccount)
	}

	wg.Wait()
	close(results)
	return len(trackingsWithAccount), m.processResults(results)
}

func (m *Monitor) maybeNotify(ctx context.Context, tracking data.TrackingAndAccount) error {
	var (
		expires       = tracking.Expires
		notifyUpfront = time.Hour * 24 * time.Duration(tracking.NotifyUpfront)
	)
	if tracking.Status != data.StatusHealthy && tracking.Status != data.StatusExpires {
		fmt.Printf("NOTIFY STATUS => %s => %s\n", tracking.DomainName, tracking.Status)
	} else if time.Until(expires) <= notifyUpfront {
		fmt.Printf("NOTIFY EXPIRES SOON => %s\n", tracking.DomainName)
	}
	return nil
}

func (m *Monitor) processResults(resultsChan chan data.DomainTracking) error {
	var (
		trackings = make([]data.DomainTracking, len(resultsChan))
		i         int
	)
	for result := range resultsChan {
		trackings[i] = result
		i++
	}
	return db.Store.Domain.UpdateAllTrackings(trackings)
}

func (m *Monitor) Start() {
	t := time.NewTicker(m.interval)
	if _, err := m.poll(); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-t.C:
			start := time.Now()
			logger.Log("msg", "new poll")
			count, err := m.poll()
			if err != nil {
				logger.Log("error", "monitor poll error", "err", err)
			}
			logger.Log("msg", "poll complete", "count", count, "took", time.Since(start))
		case <-m.quitch:
			logger.Log("msg", "monitor quitting...", "lastPoll", m.lastPoll)
			return
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db.Init()
	logger.Init()

	m := NewMonitor(time.Second * 10)
	m.Start()
}
