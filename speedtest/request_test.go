package speedtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

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

// func TestUploadTestContext(t *testing.T) {
// 	latency, _ := time.ParseDuration("5ms")
// 	server := Server{
// 		URL:     "http://fake.com/upload.php",
// 		Latency: latency,
// 	}

// 	err := server.uploadTestContext(
// 		context.Background(),
// 		false,
// 		mockWarmUp,
// 		mockRequest,
// 	)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	if server.ULSpeed < 2400 || 2600 < server.ULSpeed {
// 		t.Errorf("got unexpected server.ULSpeed '%v', expected between 2400 and 2600", server.ULSpeed)
// 	}
// }

func mockWarmUp(ctx context.Context, client *resty.Client, dlURL string, w int) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}

func mockRequest(ctx context.Context, client *resty.Client, dlURL string, w int) error {
	fmt.Sprintln(w)
	time.Sleep(500 * time.Millisecond)
	return nil
}
