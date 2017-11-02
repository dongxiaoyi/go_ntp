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

	// 监听地址和端口，地址传空表示监听所有地址。
	ntps := ntp.NewNTPS("", port)

	// 启动服务，会创建协程在后台运行。
	ntps.Start()
	for {
		time.Sleep(60 * time.Second)
	}

	// 停止服务
	ntps.Stop()
}

// 客户端demo
func Client(addr string) {

	// 建立连接
	ntpc := ntp.NewNTPC(addr)

	for {
		// 延时一段时间
		time.Sleep(5 * time.Second)

		// 创建10个同步时间结果
		resultAry := make([]ntp.Result, 10)

		// 同步执行10次，并且得到相应的结果
		for i, _ := range resultAry {
			rsp, err := ntpc.Sync(1)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			resultAry[i] = rsp
		}

		// 网络时延统计
		var result ntp.Result
		var count int64
		for _, v := range resultAry {
			if v.NetDelay.NanoSecond > 0 {
				result.NetDelay.Add(v.NetDelay)
				result.Offset.Add(v.Offset)
				count++
			}
		}

		// 同步成功次数小于50%=5/10，则放弃本次同步的数据，并且重试；
		if count < 5 {
			log.Println("sync time from ntp service failed!")
			continue
		}

		// 计算网络时延、时间偏移的平均值
		result.NetDelay.Div(count)
		result.Offset.Div(count)

		// 将网络时延和时间偏移打印到控制台。
		log.Printf("OffSet   %.3f ms \r\n", float64(result.Offset.NanoSecond)/float64(time.Millisecond))
		log.Printf("NetDelay %.3f ms \r\n", float64(result.NetDelay.NanoSecond)/float64(time.Millisecond))

		// 设置时间到host os，对于大于1s的时间，直接设置。
		// 小于1s并且大于50ms的时间偏移，采用逐渐逼近方式设置时间。
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

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}
