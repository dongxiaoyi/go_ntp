package ntp

import (
	"errors"
	"net"
	"time"
)

type NTPC struct {
	ServerAddr string        // 时间服务器地址
	RequestId  uint64        // 请求报文的序列号（用于校验）
	TimeOut    time.Duration // 通信超时时间
}

type Result struct {
	Offset   TimeStamp // 本地时间与服务器时间的差，负数表示本地时间快于服务器时间，正数表示本地时间慢于服务器时间；
	NetDelay TimeStamp // 本次请求的网络时延
}

// 申请一个NTPC客户端对象，并且初始化服务端地址
func NewNTPC(addr string, timeout time.Duration) *NTPC {
	var ntpc = NTPC{ServerAddr: addr, TimeOut: timeout}
	ntpc.RequestId = uint64(time.Now().Nanosecond())
	return &ntpc
}

// 发起时间同步请求
// 返回Result结果包括网络传输时延、时间偏移
// 发起时间同步请求
// 返回Result结果包括网络传输时延、时间偏移
func (n *NTPC) SyncBatch(times int) ([]Result, error) {
	rsp := make([]Result, 0)
	for i := 0; i < times; {
		res, err := n.SyncOnce()
		if err != nil {
			return rsp, err
			//continue
		}
		rsp = append(rsp, res)
		i++
	}
	return rsp, nil
}

// 发起时间同步请求
// 返回Result结果包括网络传输时延、时间偏移
func (n *NTPC) SyncOnce() (rsp Result, err error) {

	var buf [4096]byte
	var req Packet

	// 创建udp协议的socket服务
	socket, err := net.Dial("udp", n.ServerAddr)
	if err != nil {
		return
	}

	defer socket.Close()

	n.RequestId++

	// 初始化请求报文内容
	req.Version = 100
	req.RequestId = n.RequestId

	// 设置 read/write 超时时间
	err = socket.SetDeadline(time.Now().Add(n.TimeOut))
	if err != nil {
		return
	}

	req.T1 = TimeToTimeStamp(time.Now())

	// 序列化请求报文
	newbuf, err := CodePacket(req)
	if err != nil {
		return
	}

	// 发送到服务端
	_, err = socket.Write(newbuf)
	if err != nil {
		return
	}

	// 获取服务端应答报文
	cnt, err := socket.Read(buf[0:])
	if err != nil {
		return
	}

	// 校验应答报文大小
	if cnt != DEFAULT_PACKET_SIZE {
		err = errors.New("recv a packet not recognized")
		return
	}

	// 反序列化报文
	req, err = DecodePacket(buf[0:cnt])
	if err != nil {
		return
	}

	// 校验请求的序号是否一致
	if req.RequestId != n.RequestId {
		err = errors.New("recv a bad packet ")
		return
	}

	req.T4 = TimeToTimeStamp(time.Now())

	// 参考ntp的网络校时，计算出本地与服务器的时差
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
	t2.Add(t3)
	rsp.Offset = t2.Div(2)

	return
}

// 计算多个结果的平均值
func ResultAverage(arys []Result) Result {

	// 网络时延统计
	var res Result

	for _, v := range arys {
		if v.NetDelay.NanoSecond > 0 {
			res.NetDelay.Add(v.NetDelay)
			res.Offset.Add(v.Offset)
		}
	}

	// 计算网络时延、时间偏移的平均值
	res.NetDelay.Div(len(arys))
	res.Offset.Div(len(arys))

	return res
}
