package ntp

import (
	"log"
	"net"
	"time"
)

const (
	NTPC_DEFAULT_TIMEOUT    = 10
	NTPC_DEFAULT_RETRYTIMES = 100
)

type NTPC_Notify struct {
	NetWorkDelay   int // Microsecond
	TimeOffsetSec  int // Second
	TimeOffsetNsec int // Nanosecond
}

type NTPC struct {
	addr string

	TIMEOUT    int // wait resp timeout
	RETRYTIMES int // retry times

	gettime func() time.Time
	notify  func(NTPC_Notify)
}

func NewNTPC(ip, port string) *NTPC {

	var ntpc = NTPC{addr: ip + ":" + port}

	ntpc.TIMEOUT = NTPC_DEFAULT_TIMEOUT
	ntpc.RETRYTIMES = NTPC_DEFAULT_RETRYTIMES

	return &ntpc
}

func (n *NTPC) Config(timeout, retrytimes int) {
	n.TIMEOUT = timeout
	n.RETRYTIMES = retrytimes
}

func (n *NTPC) RegHandler(gettime func() time.Time, notify func(NTPC_Notify)) {
	n.gettime = gettime
	n.notify = notify
}

var requestid uint64

func (n *NTPC) Sync() error {

	var buf [4096]byte
	var req Packet

	socket, err := net.Dial("udp", n.addr)
	if err != nil {
		return err
	}

	defer socket.Close()

	for {

		time.Sleep(5 * time.Second)

		req.Version = 100
		req.RequestId = requestid

		requestid++

		req.T1 = GetTimeStamp()

		newbuf, err := CodePacket(req)
		if err != nil {
			return err
		}

		n, err := socket.Write(newbuf)
		if err != nil {
			return err
		}

		n, err = socket.Read(buf[0:])
		if err != nil {
			return err
		}

		if n != DEFAULT_PACKET_SIZE {
			log.Println("recv a packet not recognized ", buf[0:n])
			continue
		}

		req, err = DecodePacket(buf[0:n])
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if req.RequestId != requestid-1 {
			log.Println("recv a old packet ", req)
			continue
		}

		req.T4 = GetTimeStamp()

		log.Println(req)

		calcDiffTime(req)
	}

	return nil
}

func calcDiffTime(req Packet) {

	var offset, netdly TimeStamp
	var t1, t2, t3, t4 TimeStamp

	t1 = req.T1 // T1 客户端发送请求的时间
	t2 = req.T2 // T2 服务器接收请求的时间
	t3 = req.T3 // T3 服务器答复时间
	t4 = req.T4 // T4 客户端接收答复时间

	// 计算得出网络时延
	t2.Sub(t1)
	t4.Sub(t3)
	netdly = t2.Add(t4)

	t1 = req.T1 // T1 客户端发送请求的时间
	t2 = req.T2 // T2 服务器接收请求的时间
	t3 = req.T3 // T3 服务器答复时间
	t4 = req.T4 // T4 客户端接收答复时间

	// 计算本地与服务器时延
	t2.Sub(t1)
	t3.Sub(t4)
	offset = t2.Add(t3)

	log.Println("NetDelay: ", netdly.Sec, netdly.Nsec)
	log.Println("Offset: ", offset.Sec, offset.Nsec)

	return
}
