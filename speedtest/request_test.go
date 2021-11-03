package speedtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestPingTestContext(t *testing.T) {
	latency, _ := time.ParseDuration("1s")
	server := Server{
		URL:     "http://fake.com/upload.php",
		Latency: latency,
	}

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `test=test`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "http://fake.com/latency.txt", fakeResponder(200, resp, "text/plain"))

	err := server.pingTestContext(
		context.Background(),
		client,
	)
	assert.NoError(t, err, "unexpected error %v", err)
	assert.Less(t, server.Latency.Milliseconds(), latency.Milliseconds(), "got unexpected server.Latency '%v', expected greater than 0", server.Latency)
}

func TestPingTestContextWithStatus404(t *testing.T) {
	server := Server{
		URL: "http://fake.com/upload.php",
	}

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `test=test`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "http://fake.com/latency.txt", fakeResponder(404, resp, "text/plain"))

	err := server.pingTestContext(
		context.Background(),
		client,
	)
	assert.Error(t, err, "should expect error")
	assert.Equal(t, "unexpected status code 404 while pinging http://fake.com/latency.txt", err.Error(), "unexpected error %v", err)
}

func TestDownloadTestContext(t *testing.T) {
	latency, _ := time.ParseDuration("10ms")
	server := Server{
		URL:     "http://fake.com/upload.php",
		Latency: latency,
	}

	// Create a Resty Client
	client := resty.New()

	err := server.downloadTestContext(
		context.Background(),
		client,
		mockWarmUp,
		mockRequest,
	)
	assert.NoError(t, err, "unexpected error %v", err)
	assert.GreaterOrEqual(t, server.DLSpeed, 6300.0, "got unexpected server.DLSpeed '%v', expected between 6300 and 6600", server.DLSpeed)
	assert.LessOrEqual(t, server.DLSpeed, 6600.0, "got unexpected server.DLSpeed '%v', expected between 6300 and 6600", server.DLSpeed)
}

func TestDownloadTestContextWithStatus404(t *testing.T) {
	latency, _ := time.ParseDuration("10ms")
	server := Server{
		URL:     "http://fake.com/upload.php",
		Latency: latency,
	}

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `fake-image`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "http://fake.com/random750x750.jpg", fakeResponder(404, resp, "image/jpeg"))

	err := server.downloadTestContext(
		context.Background(),
		client,
		downloadRequest,
		downloadRequest,
	)
	assert.Error(t, err, "should expect error")
}

func TestUploadTestContext(t *testing.T) {
	latency, _ := time.ParseDuration("5ms")
	server := Server{
		URL:     "http://fake.com/upload.php",
		Latency: latency,
	}

	// Create a Resty Client
	client := resty.New()

	err := server.uploadTestContext(
		context.Background(),
		client,
		mockWarmUp,
		mockRequest,
	)
	assert.NoError(t, err, "unexpected error %v", err)
	assert.GreaterOrEqual(t, server.ULSpeed, 2400.0, "got unexpected server.ULSpeed '%v', expected between 2400 and 2600", server.ULSpeed)
	assert.LessOrEqual(t, server.ULSpeed, 2600.0, "got unexpected server.ULSpeed '%v', expected between 2400 and 2600", server.ULSpeed)
}

func TestUploadTestContextWithStatus404(t *testing.T) {
	latency, _ := time.ParseDuration("5ms")
	server := Server{
		URL:     "http://fake.com/upload.php",
		Latency: latency,
	}

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `fake-response`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("Post", "http://fake.com/upload.php", fakeResponder(404, resp, "image/jpeg"))

	err := server.uploadTestContext(
		context.Background(),
		client,
		uploadRequest,
		uploadRequest,
	)
	assert.Error(t, err, "should expect error")
}

func mockWarmUp(ctx context.Context, client *resty.Client, dlURL string, w int) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

func mockRequest(ctx context.Context, client *resty.Client, dlURL string, w int) error {
	fmt.Sprintln(w)
	time.Sleep(500 * time.Millisecond)
	return nil
}
