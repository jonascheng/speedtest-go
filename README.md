![CI](https://github.com/jonascheng/speedtest-go/actions/workflows/ci.yaml/badge.svg)
![CD](https://github.com/jonascheng/speedtest-go/actions/workflows/cd.yaml/badge.svg)
![codecov](https://codecov.io/gh/jonascheng/speedtest-go/branch/main/graph/badge.svg)

# speedtest-go
**Command Line Interface and pure Go API to Test Internet Speed using [speedtest.net](http://www.speedtest.net/)**
You can speedtest 2x faster than [speedtest.net](http://www.speedtest.net/) with almost the same result.

Inspired by [sivel/speedtest-cli](https://github.com/sivel/speedtest-cli)

### Usage

```bash
$ speedtest --help
usage: speedtest-go [<flags>]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -l, --list               Show available speedtest.net servers.
  -s, --server=SERVER ...  Select server id to speedtest.
      --json               Output results as json
      --version            Show application version.
```

#### Test Internet Speed

Simply use `speedtest` command. The closest server is selected by default.

```bash
$ ./bin/speedtest-go
Testing From IP: 211.72.129.103, (Chunghwa Telecom) (TW) [25.0504, 121.5324]

Target Server: [18445]     1.91km
	> Taipei (Taiwan) by Chunghwa Mobile
	> http://tp1.chtm.hinet.net:8080/speedtest/upload.php
Latency: 7.523354ms
Download Test: ................
Upload Test: ................

Download: 73.30 Mbit/s
Upload: 35.26 Mbit/s
```

#### Test to Other Servers

If you want to select other server to test, you can see available server list.

```bash
$ ./bin/speedtest-go --list
Testing From IP: 211.72.129.103, (Chunghwa Telecom) (TW) [25.0504, 121.5324]
[18445]     1.91km Taipei (Taiwan) by Chunghwa Mobile
[2133]     1.91km Taipei (Taiwan) by Taiwan Fixed Network
[44603]     1.91km Taipei (Taiwan) by Taiwan Digital Streaming Co.
[45693]     1.91km Taipei (Taiwan) by PEGATRON
[13506]     3.45km Taipei (Taiwan) by TAIFO Taiwan
[14652]     3.85km 新北 (Taiwan) by 大新店
[14651]     3.85km 新北 (Taiwan) by 新唐城
[17265]     7.68km Zhonghe (TW) by FarEasTone Telecom
[24461]     8.34km Banqiao (Taiwan) by Homeplus
```

and select them by id.

```bash
$ ./bin/speedtest-go --server 18445 --server 24461
Testing From IP: 211.72.129.103, (Chunghwa Telecom) (TW) [25.0504, 121.5324]

Target Server: [18445]     1.91km
	> Taipei (Taiwan) by Chunghwa Mobile
	> http://tp1.chtm.hinet.net:8080/speedtest/upload.php
Latency: 8.655645ms
Download Test: ................
Upload Test: ........

Download: 126.69 Mbit/s
Upload: 101.33 Mbit/s

Target Server: [24461]     8.34km
	> Banqiao (Taiwan) by Homeplus
	> http://sky-speedtest.bbtv.tw:8080/speedtest/upload.php
Latency: 9.962729ms
Download Test: ............................
Upload Test: ...........

Download: 113.40 Mbit/s
Upload: 128.44 Mbit/s

Download Avg: 120.04 Mbit/s
Upload Avg: 114.88 Mbit/s
```

## Go API

```
go get github.com/jonascheng/speedtest-go
```

### API Usage
The code below finds closest available speedtest server and tests the latency, download, and upload speeds.

```go
package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jonascheng/speedtest-go/speedtest"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	user, _ := speedtest.FetchUserInfo(client)

	serverList, _ := speedtest.FetchServerList(client, user)

	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		// This is required to determin network latency.
		s.PingTest(client)
		// These two bandwidth tests can be used base upon use cases.
		// If use case requires only upload bandwidth, and then just invoke UploadTest to obtain ULSpeed.
		s.DownloadTest(client)
		s.UploadTest(client)

		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
	}
}
```

## LICENSE

[MIT](https://github.com/jonascheng/speedtest-go/blob/master/LICENSE)
