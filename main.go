package main

import (
	"fmt"
	"go_ntp/ntp"
	"time"
)

func gettimebyuser() time.Time {
	return time.Now()
}

func notifytime(nty ntp.NTPC_Notify) {
	fmt.Printf("NetWorkDelay  : %d.%03d ms \r\n", nty.NetWorkDelay/1000, nty.NetWorkDelay%1000)
	fmt.Printf("TimeOffset    : %d.%09d s.ns\r\n", nty.TimeOffsetSec, nty.TimeOffsetNsec)
}

func main() {
	ntpc := ntp.NewNTPC("time.nist.gov", "123")

	ntpc.Config(2, 30)

	ntpc.RegHandler(gettimebyuser, notifytime)

	for i := 0; i < 10000; i++ {
		err := ntpc.Sync()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	return
}
