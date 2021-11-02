package speedtest

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchServerList(t *testing.T) {
	// Create a Resty Client
	client := resty.New()

	user := User{
		IP:  "111.111.111.111",
		Lat: "35.22",
		Lon: "138.44",
		Isp: "Hello",
	}
	serverList, err := FetchServerList(client, &user)
	assert.NoError(t, err, "unexpected error")
	assert.Greater(t, len(serverList.Servers), 0, "failed to fetch server list.")
	assert.Greater(t, len(serverList.Servers[0].Country), 0, "got unexpected country name '%v'", serverList.Servers[0].Country)
}

func TestFetchServerListWithFakeResponse(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `<settings>
	<servers>
	<server url="http://fake.com:8080/speedtest/upload.php" lat="35.22" lon="138.44" name="新北" country="Taiwan" cc="TW" sponsor="大新店" id="14652" host="fake.com:8080"/>
	</servers>
	</settings>`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", speedTestServersUrl, fakeResponder(200, resp, "application/xml"))

	user := User{
		IP:      "111.111.111.111",
		Lat:     "35.22",
		Lon:     "138.44",
		Isp:     "Hello",
		Country: "US",
	}
	serverList, err := FetchServerList(client, &user)
	assert.NoError(t, err, "unexpected error")
	assert.Greater(t, len(serverList.Servers), 0, "failed to fetch server list.")
	assert.Equal(t, "http://fake.com:8080/speedtest/upload.php", serverList.Servers[0].URL)
	assert.Equal(t, "新北", serverList.Servers[0].Name)
	assert.Equal(t, "Taiwan", serverList.Servers[0].Country)
	assert.Equal(t, "大新店", serverList.Servers[0].Sponsor)
	assert.Equal(t, "14652", serverList.Servers[0].ID)
	assert.Equal(t, "fake.com:8080", serverList.Servers[0].Host)
	d := distance(35.22, 138.44, 35.22, 138.44)
	assert.Equal(t, d, serverList.Servers[0].Distance, "the distance should be the same")
}

func TestFetchServerListWithEmptyResponse(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `<settings></settings>`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", speedTestServersUrl, fakeResponder(200, resp, "application/xml"))

	user := User{
		IP:      "111.111.111.111",
		Lat:     "35.22",
		Lon:     "138.44",
		Isp:     "Hello",
		Country: "US",
	}
	serverList, err := FetchServerList(client, &user)
	assert.Error(t, err, "should expect error")
	assert.Equal(t, "unable to retrieve server list", err.Error(), "unexpected error")
	assert.Equal(t, ServerList{}, serverList)
}

func TestDistance(t *testing.T) {
	d := distance(0.0, 0.0, 1.0, 1.0)
	assert.GreaterOrEqual(t, d, 157.0, "got: %v, expected between 157 and 158", d)
	assert.LessOrEqual(t, d, 158.0, "got: %v, expected between 157 and 158", d)

	d = distance(0.0, 180.0, 0.0, -180.0)
	assert.Equal(t, d, 0.0, "got: %v, expected 0", d)

	d1 := distance(100.0, 100.0, 100.0, 101.0)
	d2 := distance(100.0, 100.0, 100.0, 99.0)
	assert.Equal(t, d1, d2, "%v and %v should be save value", d1, d2)

	d = distance(35.0, 140.0, -40.0, -140.0)
	assert.GreaterOrEqual(t, d, 11000.0, "got: %v, expected between 11000 and 12000", d)
	assert.LessOrEqual(t, d, 12000.0, "got: %v, expected between 11000 and 12000", d)
}

func TestFindServer(t *testing.T) {
	servers := []*Server{
		&Server{
			ID: "1",
		},
		&Server{
			ID: "2",
		},
		&Server{
			ID: "3",
		},
	}
	serverList := ServerList{Servers: servers}

	serverID := []int{}
	s, err := serverList.FindServer(serverID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(s) != 1 {
		t.Errorf("Unexpected server length. got: %v, expected: 1", len(s))
	}
	if s[0].ID != "1" {
		t.Errorf("Unexpected server ID. got: %v, expected: '1'", s[0].ID)
	}

	serverID = []int{2}
	s, err = serverList.FindServer(serverID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(s) != 1 {
		t.Errorf("Unexpected server length. got: %v, expected: 1", len(s))
	}
	if s[0].ID != "2" {
		t.Errorf("Unexpected server ID. got: %v, expected: '2'", s[0].ID)
	}

	serverID = []int{3, 1}
	s, err = serverList.FindServer(serverID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(s) != 2 {
		t.Errorf("Unexpected server length. got: %v, expected: 2", len(s))
	}
	if s[0].ID != "3" {
		t.Errorf("Unexpected server ID. got: %v, expected: '3'", s[0].ID)
	}
	if s[1].ID != "1" {
		t.Errorf("Unexpected server ID. got: %v, expected: '1'", s[0].ID)
	}
}
