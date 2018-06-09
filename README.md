# go_ntp
A simple time synchronization library based on NTP principle

## Features

- Time synchronization [service](https://github.com/lixiangyun/go_ntp/blob/master/ntp/server.go) & [client](https://github.com/lixiangyun/go_ntp/blob/master/ntp/client.go)
- Usage method, Please to see [main.go](https://github.com/lixiangyun/go_ntp/blob/master/main.go)

## Example

### As a service example code.

```
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lixiangyun/go_ntp"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: < IP:PORT >")
		return
	}

	// Listens to addresses and ports
	ntps := ntp.NewNTPS(args[1])

	// Start the service, then coroutine process is created in the background.
	ntps.Start()
	for {
		time.Sleep(60 * time.Second)
	}

	// Stop the service.
	ntps.Stop()
}

```

#### Build

```
go build main.go
```

#### Run

```
./ntpserver :123
```

### As a client example code.

```
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lixiangyun/go_ntp"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: < IP:PORT >")
		return
	}

	/* Connect to the ntp server. */
	ntpc := ntp.NewNTPC(args[1], 1*time.Second)

	for {
		// run Delay
		time.Sleep(2 * time.Second)

		/* Initiating synchronization time, waiting for 10 results. */
		resultAry := ntpc.SyncBatch(10)

		result := ntp.ResultAverage(resultAry)

		/* Print the network delay and time offset to the console. */
		log.Printf("Network offset %d us \r\n",
			result.Offset.NanoSecond/int64(time.Microsecond))

		log.Printf("Network delay  %d us \r\n",
			result.NetDelay.NanoSecond/int64(time.Microsecond))
	}
}

```

#### Build

```
go build main.go
```

#### Run

```
./ntpclient :123
```

#### Output

```
2018/06/09 20:55:04 Network offset -49 us
2018/06/09 20:55:04 Network delay  99 us
2018/06/09 20:55:06 Network offset 49 us
2018/06/09 20:55:06 Network delay  298 us
```
