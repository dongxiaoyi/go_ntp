package main

import (
	"fmt"
	"go_ntp/ntp"
	"os"
	"time"
)

// 服务端demo
func Server() {
	ntps := ntp.NewNTPS("", "1234")
	ntps.Start()
	for {
		time.Sleep(1 * time.Second)
	}
	ntps.Stop()
}

// 客户端demo
func Client() {
	ntpc := ntp.NewNTPC("192.168.0.101", "1234")

	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Second)

		result, err := ntpc.Sync(10)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Printf("offset %.3f ms \r\n", float64(result.Offset.NanoSecond)/float64(time.Millisecond))
		fmt.Printf("netdelay %.3f ms \r\n", float64(result.NetDelay.NanoSecond)/float64(time.Millisecond))

		if result.Offset.Abs() > int64(time.Second) {
			now := ntp.TimeStampToTime(result.Offset, time.Now())
			err = ntp.SetTimeToOs(now)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: <-s/-c>")
		return
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
