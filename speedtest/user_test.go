package speedtest

import (
	"strconv"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchUserInfo(t *testing.T) {
	// Create a Resty Client
	client := resty.New()

	user, err := FetchUserInfo(client)
	if err != nil {
		t.Errorf(err.Error())
	}
	// IP
	if len(user.IP) < 7 || len(user.IP) > 15 {
		t.Errorf("Invalid IP length. got: %v;", user.IP)
	}
	if strings.Count(user.IP, ".") != 3 {
		t.Errorf("Invalid IP format. got: %v", user.IP)
	}
	// Lat
	lat, err := strconv.ParseFloat(user.Lat, 64)
	if err != nil {
		t.Errorf(err.Error())
	}
	if lat < -90 || 90 < lat {
		t.Errorf("Invalid Latitude. got: %v, expected between -90 and 90", user.Lat)
	}
	// Lon
	lon, err := strconv.ParseFloat(user.Lon, 64)
	if err != nil {
		t.Errorf(err.Error())
	}
	if lon < -180 || 180 < lon {
		t.Errorf("Invalid Latitude. got: %v, expected between -90 and 90", user.Lon)
	}
	// Isp
	if len(user.Isp) == 0 {
		t.Errorf("Invalid Iso. got: %v;", user.Isp)
	}
}

func fakeResponder(s int, c string, ct string) httpmock.Responder {
	resp := httpmock.NewStringResponse(s, c)
	resp.Header.Set("Content-Type", ct)

	return httpmock.ResponderFromResponse(resp)
}

func TestFetchUserInfoWithFakeResponse(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	// Create a Resty Client
	client := resty.New()

	// fake response
	resp := `<settings>
	<client ip="211.72.129.103" lat="25.0504" lon="121.5324" isp="Chunghwa Telecom" isprating="3.7" rating="0" ispdlavg="0" ispulavg="0" loggedin="0" country="TW"/>
	<server-config threadcount="4" ignoreids="" notonmap="" forcepingid="" preferredserverid=""/>
	<licensekey>f7a45ced624d3a70-1df5b7cd427370f7-b91ee21d6cb22d7b</licensekey>
	<customer>speedtest</customer>
	<odometer start="19601573884" rate="12"/>
	<times dl1="5000000" dl2="35000000" dl3="800000000" ul1="1000000" ul2="8000000" ul3="35000000"/>
	<download testlength="10" initialtest="250K" mintestsize="250K" threadsperurl="4"/>
	<upload testlength="10" ratio="5" initialtest="0" mintestsize="32K" threads="2" maxchunksize="512K" maxchunkcount="50" threadsperurl="4"/>
	<latency testlength="10" waittime="50" timeout="20"/>
	<socket-download testlength="15" initialthreads="4" minthreads="4" maxthreads="32" threadratio="750K" maxsamplesize="5000000" minsamplesize="32000" startsamplesize="1000000" startbuffersize="1" bufferlength="5000" packetlength="1000" readbuffer="65536"/>
	<socket-upload testlength="15" initialthreads="dyn:tcpulthreads" minthreads="dyn:tcpulthreads" maxthreads="32" threadratio="750K" maxsamplesize="1000000" minsamplesize="32000" startsamplesize="100000" startbuffersize="2" bufferlength="1000" packetlength="1000" disabled="false"/>
	<socket-latency testlength="10" waittime="50" timeout="20"/>
	<translation lang="xml"> </translation>
	</settings>`

	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "https://www.speedtest.net/speedtest-config.php", fakeResponder(200, resp, "application/xml"))

	user, err := FetchUserInfo(client)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, "211.72.129.103", user.IP)
	assert.Equal(t, "25.0504", user.Lat)
	assert.Equal(t, "121.5324", user.Lon)
	assert.Equal(t, "Chunghwa Telecom", user.Isp)
	assert.Equal(t, "TW", user.Country)
}
