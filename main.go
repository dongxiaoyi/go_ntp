package main

import (
	"fmt"
	"go_ntp/ntp"
	"os"
	"time"
)

func Server() {
	ntps := ntp.NewNTPS("", "1234")
	ntps.Start()
	for {
		time.Sleep(1 * time.Second)
	}
	ntps.Stop()
}

func Client() {
	ntpc := ntp.NewNTPC("192.168.0.107", "1234")

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)

		result, err := ntpc.Sync(10)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Printf("offset %d.%09d sec.nanosec \r\n", result.Offset.Sec, result.Offset.Nsec)
		fmt.Printf("netdelay %.3f ms \r\n", float32(result.NetDelay.Nsec)/float32(time.Millisecond))
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: <-s/-c>")
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
