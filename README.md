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
      --saving-mode        Using less memory (≒10MB), though low accuracy (especially > 30Mbps).
      --json               Output results as json
      --version            Show application version.
```

#### Test Internet Speed

Simply use `speedtest` command. The closest server is selected by default.

```bash
$ speedtest
Testing From IP: 124.27.199.165 (Fujitsu) [34.9769, 138.3831]

Target Server: [6691]     9.03km Shizuoka (Japan) by sudosan
latency: 39.436061ms
Download Test: ................
Upload Test: ................

Download: 73.30 Mbit/s
Upload: 35.26 Mbit/s
```

#### Test to Other Servers

If you want to select other server to test, you can see available server list.

```bash
$ speedtest --list
Testing From IP: 124.27.199.165 (Fujitsu) [34.9769, 138.3831]
[6691]     9.03km Shizuoka (Japan) by sudosan
[6087]   120.55km Fussa-shi (Japan) by Allied Telesis Capital Corporation
[6508]   125.44km Yokohama (Japan) by at2wn
[6424]   148.23km Tokyo (Japan) by Cordeos Corp.
[6492]   153.06km Sumida (Japan) by denpa893
[7139]   192.63km Tsukuba (Japan) by SoftEther Corporation
[6368]   194.83km Maibara (Japan) by gatolabo
[6463]   220.39km Kusatsu (Japan) by j416dy
[6766]   232.54km Nomi (Japan) by JAIST(ino-lab)
[6476]   265.10km Osaka (Japan) by rxy (individual)
[6477]   268.94km Sakai (Japan) by satoweb
...
```

and select them by id.

```bash
$ speedtest --server 6691 --server 6087
Testing From IP: 124.27.199.165 (Fujitsu) [34.9769, 138.3831]

Target Server: [6691]     9.03km Shizuoka (Japan) by sudosan
Latency: 23.612861ms
Download Test: ................
Upload Test: ........

Target Server: [6087]   120.55km Fussa-shi (Japan) by Allied Telesis Capital Corporation
Latency: 38.694699ms
Download Test: ................
Upload Test: ................

[6691] Download: 65.82 Mbit/s, Upload: 27.00 Mbit/s
[6087] Download: 72.24 Mbit/s, Upload: 29.56 Mbit/s
Download Avg: 69.03 Mbit/s
Upload Avg: 28.28 Mbit/s
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
	"github.com/jonascheng/speedtest-go/speedtest"
)

func main() {
	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
	}
}
```

## LICENSE

[MIT](https://github.com/jonascheng/speedtest-go/blob/master/LICENSE)
