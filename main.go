package main

import (
	"fmt"
	"go_ntp/ntp"
	"os"
	"time"
)

// 服务端demo
func Server() {
	ntps := ntp.NewNTPS("dev0", "3210")
	ntps.Start()
	for {
		time.Sleep(1 * time.Second)
	}
	ntps.Stop()
}

// 客户端demo
func Client() {

	ntpc := ntp.NewNTPC("dev0", "3210")

	for {
		time.Sleep(1 * time.Second)

		resultAry := make([]ntp.Result, 10)

		for i, _ := range resultAry {
			rsp, err := ntpc.Sync(1)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			resultAry[i] = rsp
		}

		var result ntp.Result
		var count int64
		for _, v := range resultAry {
			if v.NetDelay.NanoSecond > 0 {
				result.NetDelay.Add(v.NetDelay)
				result.Offset.Add(v.Offset)
				count++
			}
		}

		if count < 5 {
			fmt.Println("sync time from ntp service failed!")
			continue
		}

		result.NetDelay.Div(count)
		result.Offset.Div(count)

		fmt.Printf("OffSet   %.3f ms \r\n", float64(result.Offset.NanoSecond)/float64(time.Millisecond))
		fmt.Printf("NetDelay %.3f ms \r\n", float64(result.NetDelay.NanoSecond)/float64(time.Millisecond))

		if result.Offset.Abs() > int64(time.Second) {
			now := ntp.TimeStampToTime(result.Offset, time.Now())
			fmt.Println(result.Offset)

			err := ntp.SetTimeToOs(now)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		} else if result.Offset.Abs() > 50*int64(time.Millisecond) {

			now := ntp.TimeStampToTime(result.Offset.Div(4), time.Now())
			fmt.Println(result.Offset)

			err := ntp.SetTimeToOs(now)
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
