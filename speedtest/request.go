package speedtest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
)

type downloadFunc func(context.Context, *resty.Client, string, int) error
type uploadFunc func(context.Context, *resty.Client, string, int) error

var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000} //kB

// DownloadTest executes the test to measure download speed
func (s *Server) DownloadTest(client *resty.Client) error {
	return s.downloadTestContext(context.Background(), client, downloadRequest, downloadRequest)
}

func (s *Server) downloadTestContext(
	ctx context.Context,
	client *resty.Client,
	dlWarmUp downloadFunc,
	downloadRequest downloadFunc,
) error {
	dlURL := strings.Split(s.URL, "/upload.php")[0]
	eg := errgroup.Group{}

	// Warming up
	wuWeight := 2
	sTime := time.Now()
	for i := 0; i < 2; i++ {
		eg.Go(func() error {
			return dlWarmUp(ctx, client, dlURL, wuWeight)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	fTime := time.Now()
	// 1.125MB for each request (750 * 750 * 2  / 1000 / 1000)
	reqMB := float64(dlSizes[wuWeight]) * float64(dlSizes[wuWeight]) * 2.0 / 1000.0 / 1000.0
	// Calculate speed in Mbps
	wuSpeed := float64(reqMB) * 8.0 * 2.0 / fTime.Sub(sTime.Add(s.Latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	switch {
	case wuSpeed > 50.0:
		workload = 32
		weight = 6
	case wuSpeed > 10.0:
		workload = 16
		weight = 4
	case wuSpeed > 4.0:
		workload = 8
		weight = 4
	case wuSpeed > 2.5:
		workload = 4
		weight = 4
	default:
		skip = true
	}

	// Main speedtest
	dlSpeed := wuSpeed
	if !skip {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			eg.Go(func() error {
				return downloadRequest(ctx, client, dlURL, weight)
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		fTime = time.Now()

		reqMB := float64(dlSizes[weight]) * float64(dlSizes[weight]) * 2.0 / 1000.0 / 1000.0
		dlSpeed = float64(reqMB) * 8 * float64(workload) / fTime.Sub(sTime.Add(s.Latency)).Seconds()
	}

	s.DLSpeed = dlSpeed
	return nil
}

// UploadTest executes the test to measure upload speed
func (s *Server) UploadTest(client *resty.Client) error {
	return s.uploadTestContext(context.Background(), client, uploadRequest, uploadRequest)
}

func (s *Server) uploadTestContext(
	ctx context.Context,
	client *resty.Client,
	ulWarmUp uploadFunc,
	uploadRequest uploadFunc,
) error {
	// Warm up
	wuWeight := 4
	sTime := time.Now()
	eg := errgroup.Group{}
	for i := 0; i < 2; i++ {
		eg.Go(func() error {
			return ulWarmUp(ctx, client, s.URL, wuWeight)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	fTime := time.Now()
	// 1.0 MB for each request
	reqMB := float64(ulSizes[wuWeight]) / 1000.0
	wuSpeed := reqMB * 8 * 2.0 / fTime.Sub(sTime.Add(s.Latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	switch {
	case wuSpeed > 50.0:
		workload = 40
		weight = 9
	case wuSpeed > 10.0:
		workload = 16
		weight = 9
	case wuSpeed > 4.0:
		workload = 8
		weight = 9
	case wuSpeed > 2.5:
		workload = 4
		weight = 5
	default:
		skip = true
	}

	// Main speedtest
	ulSpeed := wuSpeed
	if !skip {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			eg.Go(func() error {
				return uploadRequest(ctx, client, s.URL, weight)
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		fTime = time.Now()

		reqMB := float64(ulSizes[weight]) / 1000.0
		ulSpeed = reqMB * 8 * float64(workload) / fTime.Sub(sTime.Add(s.Latency)).Seconds()
	}

	s.ULSpeed = ulSpeed

	return nil
}

func downloadRequest(ctx context.Context, client *resty.Client, dlURL string, w int) error {
	size := dlSizes[w]
	xdlURL := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"

	resp, err := client.R().
		SetContext(ctx).
		Get(xdlURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code %v while downloading from %v", resp.StatusCode(), xdlURL)
	}

	return err
}

func uploadRequest(ctx context.Context, client *resty.Client, ulURL string, w int) error {
	size := ulSizes[w]
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size*100-51))

	resp, err := client.R().
		SetContext(ctx).
		SetBody(v.Encode()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post(ulURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code %v while uploading to %v", resp.StatusCode(), ulURL)
	}

	return err
}

// PingTest executes test to measure latency
func (s *Server) PingTest(client *resty.Client) error {
	return s.pingTestContext(context.Background(), client)
}

// pingTestContext executes test to measure latency, observing the given context.
func (s *Server) pingTestContext(ctx context.Context, client *resty.Client) error {
	pingURL := strings.Split(s.URL, "/upload.php")[0] + "/latency.txt"

	l := time.Duration(10000000000) // 10sec
	for i := 0; i < 3; i++ {
		sTime := time.Now()

		resp, err := client.R().
			SetContext(ctx).
			Get(pingURL)

		if err != nil {
			return err
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code %v while pinging %v", resp.StatusCode(), pingURL)
		}

		fTime := time.Now()
		if fTime.Sub(sTime) < l {
			l = fTime.Sub(sTime)
		}
	}

	// divide by 2 due to round trip time per request
	s.Latency = time.Duration(int64(l.Nanoseconds() / 2))

	return nil
}
