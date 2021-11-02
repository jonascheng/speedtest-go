package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/jonascheng/speedtest-go/speedtest"
	// A Go (golang) command line and flag parser
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	showList   = kingpin.Flag("list", "Show available speedtest.net servers.").Short('l').Bool()
	serverIds  = kingpin.Flag("server", "Select server id to speedtest.").Short('s').Ints()
	jsonOutput = kingpin.Flag("json", "Output results in json format").Bool()
)

type fullOutput struct {
	UserInfo *speedtest.User   `json:"user_info"`
	Servers  speedtest.Servers `json:"servers"`
}

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	// Create a Resty Client
	client := resty.New()
	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second)

	user, err := speedtest.FetchUserInfo(client)
	if err != nil {
		fmt.Println("Warning: Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
	}
	if !*jsonOutput {
		showUser(user)
	}

	serverList, err := speedtest.FetchServerList(client, user)
	checkError(err)
	if *showList {
		showServerList(serverList)
		return
	}

	targets, err := serverList.FindServer(*serverIds)
	checkError(err)

	startTest(client, targets, *jsonOutput)

	if *jsonOutput {
		jsonBytes, err := json.MarshalIndent(
			fullOutput{
				UserInfo: user,
				Servers:  targets,
			},
			"",
			"  ",
		)
		checkError(err)

		fmt.Println(string(jsonBytes))
	}
}

func startTest(client *resty.Client, servers speedtest.Servers, jsonOutput bool) {
	for _, s := range servers {
		if !jsonOutput {
			showServer(s)
		}

		err := s.PingTest(client)
		checkError(err)

		if jsonOutput {
			err := s.DownloadTest(client)
			checkError(err)

			err = s.UploadTest(client)
			checkError(err)

			continue
		}

		showLatencyResult(s)

		err = testDownload(s, client)
		checkError(err)
		err = testUpload(s, client)
		checkError(err)

		showServerResult(s)
	}

	if !jsonOutput && len(servers) > 1 {
		showAverageServerResult(servers)
	}
}

func testDownload(server *speedtest.Server, client *resty.Client) error {
	quit := make(chan bool)
	fmt.Printf("Download Test: ")
	go dots(quit)
	err := server.DownloadTest(client)
	quit <- true
	if err != nil {
		return err
	}
	fmt.Println()
	return err
}

func testUpload(server *speedtest.Server, client *resty.Client) error {
	quit := make(chan bool)
	fmt.Printf("Upload Test: ")
	go dots(quit)
	err := server.UploadTest(client)
	quit <- true
	if err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func dots(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
			fmt.Print(".")
		}
	}
}

func showUser(user *speedtest.User) {
	if user.IP != "" {
		fmt.Printf("Testing From IP: %s\n", user.String())
	}
}

func showServerList(serverList speedtest.ServerList) {
	for _, s := range serverList.Servers {
		fmt.Printf("[%4s] %8.2fkm ", s.ID, s.Distance)
		fmt.Printf(s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
	}
}

func showServer(s *speedtest.Server) {
	fmt.Printf(" \n")
	fmt.Printf("Target Server: [%4s] %8.2fkm\n", s.ID, s.Distance)
	fmt.Printf("\t> " + s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
	fmt.Printf("\t> " + s.URL + "\n")
}

func showLatencyResult(server *speedtest.Server) {
	fmt.Println("Latency:", server.Latency)
}

// ShowResult : show testing result
func showServerResult(server *speedtest.Server) {
	fmt.Printf(" \n")

	fmt.Printf("Download: %5.2f Mbit/s\n", server.DLSpeed)
	fmt.Printf("Upload: %5.2f Mbit/s\n\n", server.ULSpeed)
	valid := server.CheckResultValid()
	if !valid {
		fmt.Println("Warning: Result seems to be wrong. Please speedtest again.")
	}
}

func showAverageServerResult(servers speedtest.Servers) {
	avgDL := 0.0
	avgUL := 0.0
	for _, s := range servers {
		avgDL = avgDL + s.DLSpeed
		avgUL = avgUL + s.ULSpeed
	}
	fmt.Printf("Download Avg: %5.2f Mbit/s\n", avgDL/float64(len(servers)))
	fmt.Printf("Upload Avg: %5.2f Mbit/s\n", avgUL/float64(len(servers)))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
