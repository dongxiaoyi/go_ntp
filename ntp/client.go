package ntp

import (
	"errors"
	"net"
	"time"
)

type NTPC struct {
	ServerAddr string
	RequestId  uint64
}

type Result struct {
	Offset, NetDelay TimeStamp
}

func NewNTPC(ip, port string) *NTPC {
	var ntpc = NTPC{ServerAddr: ip + ":" + port}
	ntpc.RequestId = uint64(time.Now().Nanosecond())
	return &ntpc
}

func (n *NTPC) Sync(timeout int) (rsp Result, err error) {

	var buf [4096]byte
	var req Packet

	socket, err := net.Dial("udp", n.ServerAddr)
	if err != nil {
		return
	}

	defer socket.Close()

	n.RequestId++

	req.Version = 100
	req.RequestId = n.RequestId

	err = socket.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if err != nil {
		return
	}

	req.T1 = GetTimeStamp()

	newbuf, err := CodePacket(req)
	if err != nil {
		return
	}

	_, err = socket.Write(newbuf)
	if err != nil {
		return
	}

	cnt, err := socket.Read(buf[0:])
	if err != nil {
		return
	}

	if cnt != DEFAULT_PACKET_SIZE {
		err = errors.New("recv a packet not recognized")
		return
	}

	req, err = DecodePacket(buf[0:cnt])
	if err != nil {
		return
	}

	if req.RequestId != n.RequestId {
		err = errors.New("recv a bad packet ")
		return
	}

	req.T4 = GetTimeStamp()

	rsp = calcDiffTime(req)

	return
}

func calcDiffTime(req Packet) (rsp Result) {
	var t1, t2, t3, t4 TimeStamp

	t1 = req.T1 // T1 客户端发送请求的时间
	t2 = req.T2 // T2 服务器接收请求的时间
	t3 = req.T3 // T3 服务器答复时间
	t4 = req.T4 // T4 客户端接收答复时间

	// 计算得出网络时延
	t2.Sub(t1)
	t4.Sub(t3)
	rsp.NetDelay = t2.Add(t4)

	t1 = req.T1 // T1 客户端发送请求的时间
	t2 = req.T2 // T2 服务器接收请求的时间
	t3 = req.T3 // T3 服务器答复时间
	t4 = req.T4 // T4 客户端接收答复时间

	// 计算本地与服务器时延
	t2.Sub(t1)
	t3.Sub(t4)
	rsp.Offset = t2.Add(t3)

	return
}
