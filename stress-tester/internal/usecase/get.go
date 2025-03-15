package usecase

import (
	"context"
	"fmt"
	"net/http"
	"stress-tester/internal/db"
	"stress-tester/internal/dto"
	"stress-tester/internal/entity"
	"stress-tester/internal/pool"
	"stress-tester/internal/report"
	"stress-tester/internal/stats"
	"sync"
	"time"
)

type httpGet struct {
	Client        *http.Client
	Target        string
	ReturnChannel chan *dto.Red
	NumRequests   int
}

// newHttpGet creates an httpGet object with the given http client, target, number of requests and
// return channel.
func newHttpGet(client *http.Client, target string, numRequests int, rec chan *dto.Red) *httpGet {
	return &httpGet{
		Client:        client,
		Target:        target,
		ReturnChannel: rec,
		NumRequests:   numRequests,
	}
}
// executeGet runs the http gets in a loop, stopping when the context is canceled,
// and sends the results of each get down the channel.
func (h *httpGet) executeGet(ctx context.Context, wg *sync.WaitGroup) {
	for range h.NumRequests {
		select {
		case <-ctx.Done():
			return
		default:
			go func(client *http.Client, target string, rec chan *dto.Red, wg *sync.WaitGroup) {
				r := &entity.Red{
					Target: target,
				}
				r.Get(client)
				dto := &dto.Red{Target: r.Target, SentAt: r.SentAt, ReceivedAt: r.ReceivedAt, StatusCode: r.StatusCode, Duration: r.ReceivedAt.Sub(r.SentAt)}
				rec <- dto
				wg.Done()
			}(h.Client, h.Target, h.ReturnChannel, wg)
		}
	}
}

// RoutineGet runs a number of GET requests against a target url and stores the responses in
// a database. It will run the given number of requests, but will do so in batches of
// concurrency. It will cancel any remaining work when all requests have been completed.
// It will then generate a report on the stored data and print it to the console.
func RoutineGet(target string, requests int, concurrency int) {
	start := time.Now()

	rounds := int(float64(requests) / float64(concurrency))
	extra := requests - concurrency*rounds

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rec := make(chan *dto.Red)

	database := db.NewDB(pool.GetDb(), rec)
	defer database.Close()

	go database.Store(ctx)

	wg := sync.WaitGroup{}

	for i := range rounds {
		fmt.Println("Round ", i, "Running ", concurrency, " requests for endpoint ", target)
		hg := newHttpGet(pool.GetHttpClient(), target, concurrency, rec)
		wg.Add(concurrency)
		hg.executeGet(ctx, &wg)
	}

	if extra > 0 {
		fmt.Println("Round ", rounds, "Running ", extra, " requests for endpoint ", target)
		hg := newHttpGet(pool.GetHttpClient(), target, extra, rec)
		wg.Add(extra)
		hg.executeGet(ctx, &wg)
	}

	wg.Wait()
	time.Sleep(time.Millisecond)
	cancel()

	fmt.Println("Finished ", requests, " requests for endpoint ", target, " in ", time.Since(start))
	report.ReportRed(stats.CalculateRed(database.GetAllReds()))
	report.ReportError(stats.CalculateErrors(database.GetAllReds()))
	report.ReportPercentiles(stats.CalculatePercentile(database.GetAllReds()))
	database.Close()
	time.Sleep(time.Second)
}
