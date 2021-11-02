package speedtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const speedTestConfigUrl = "https://www.speedtest.net/speedtest-config.php"

// User represents information determined about the caller by speedtest.net
type User struct {
	// <client ip="211.72.129.103" lat="25.0504" lon="121.5324" isp="Chunghwa Telecom" country="TW"/>
	IP      string `xml:"ip,attr"`
	Lat     string `xml:"lat,attr"`
	Lon     string `xml:"lon,attr"`
	Isp     string `xml:"isp,attr"`
	Country string `xml:"country,attr"`
}

// Users for decode xml
type Users struct {
	Users []User `xml:"client"`
}

// FetchUserInfo returns information about caller determined by speedtest.net
func FetchUserInfo(client *resty.Client) (*User, error) {
	return FetchUserInfoContext(context.Background(), client)
}

// FetchUserInfoContext returns information about caller determined by speedtest.net, observing the given context.
func FetchUserInfoContext(ctx context.Context, client *resty.Client) (*User, error) {
	var users Users

	_, err := client.R().
		SetContext(ctx).
		SetResult(&users).
		Get(speedTestConfigUrl)

	if err != nil {
		return nil, err
	}

	if len(users.Users) == 0 {
		return nil, errors.New("failed to fetch user information")
	}

	// Only return the first item
	return &users.Users[0], nil
}

// String representation of User
func (u *User) String() string {
	return fmt.Sprintf("%s, (%s) (%s) [%s, %s]", u.IP, u.Isp, u.Country, u.Lat, u.Lon)
}
