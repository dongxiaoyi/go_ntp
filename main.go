package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lixiangyun/go_ntp/ntp"
)

// 服务端demo
func Server(port string) {
	ntps := ntp.NewNTPS("", port)
	ntps.Start()
	for {
		time.Sleep(60 * time.Second)
	}
	ntps.Stop()
}

// 客户端demo
func Client(addr string) {
	ntpc := ntp.NewNTPC(addr)
	for {

		time.Sleep(5 * time.Second)

		resultAry := make([]ntp.Result, 10)

		for i, _ := range resultAry {
			rsp, err := ntpc.Sync(1)
			if err != nil {
				log.Println(err.Error())
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
			log.Println("sync time from ntp service failed!")
			continue
		}

		result.NetDelay.Div(count)
		result.Offset.Div(count)

		log.Printf("OffSet   %.3f ms \r\n", float64(result.Offset.NanoSecond)/float64(time.Millisecond))
		log.Printf("NetDelay %.3f ms \r\n", float64(result.NetDelay.NanoSecond)/float64(time.Millisecond))

		if result.Offset.Abs() > int64(time.Second) {
			now := ntp.TimeStampToTime(result.Offset, time.Now())
			log.Println(result.Offset)

			err := ntp.SetTimeToOs(now)
			if err != nil {
				log.Println(err.Error())
				break
			}
		} else if result.Offset.Abs() > 50*int64(time.Millisecond) {

			now := ntp.TimeStampToTime(result.Offset.Div(4), time.Now())
			log.Println(result.Offset)

			err := ntp.SetTimeToOs(now)
			if err != nil {
				log.Println(err.Error())
				break
			}
		}
	}
}

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: < -s PORT / -c IP:PORT >")
		return
	}

	file, err := os.OpenFile("runlog.txt", os.O_WRONLY, 0)
	if err != nil {
		file, err = os.Create("runlog.txt")
		if err != nil {
			fmt.Println("create file error!")
			return
		}
		fmt.Println("create log file ")
	} else {
		fileinfo, err := file.Stat()
		if err == nil {
			file.Seek(fileinfo.Size(), 0)
			fmt.Println("append log to file ", file.Name())
		}
	}

	defer file.Close()

	log.SetOutput(file)

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}
