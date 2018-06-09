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
