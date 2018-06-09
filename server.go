package ntp

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

const (
	DEFAULT_PACKET_SIZE = 6 * 8 // udp报文大小
)

// UDP报文结构信息
type Packet struct {
	Version   uint64    // 版本信息，目前暂未使用
	RequestId uint64    // 请求编号，用于校验
	T1        TimeStamp // T1 客户端发送请求的时间
	T2        TimeStamp // T2 服务器接收请求的时间
	T3        TimeStamp // T3 服务器答复时间
	T4        TimeStamp // T4 客户端接收答复时间
}

// 时间服务端结构
type NTPS struct {
	request int64
	addr    string
	conn    *net.UDPConn
	wait    sync.WaitGroup
}

// 报文序列化
func CodePacket(req Packet) ([]byte, error) {
	iobuf := new(bytes.Buffer)

	err := binary.Write(iobuf, binary.BigEndian, req)
	if err != nil {
		return nil, err
	}

	//log.Println("REQ: ", req)
	//log.Println("SEND_BUF: ", iobuf.Len(), iobuf.Bytes())

	return iobuf.Bytes(), nil
}

// 报文反序列化
func DecodePacket(buf []byte) (rsp Packet, err error) {

	iobuf := bytes.NewReader(buf)
	err = binary.Read(iobuf, binary.BigEndian, &rsp)

	//log.Println("RSP: ", rsp)
	//log.Println("RECV_BUF:", len(buf), buf)

	return
}

// 消息收发的处理协成
func msgProc(s *NTPS) {

	defer s.wait.Done()
	var buf [4096]byte

	for {
		// 监听
		n, addr, err := s.conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println("socket close.")
			return
		}

		// 获取服务器本地时间
		T2 := TimeToTimeStamp(time.Now())

		// 校验报文大小是否符合预期
		if n != DEFAULT_PACKET_SIZE {
			log.Println("recv a packet not recognized ", len(buf), buf[0:n])
			continue
		}

		// 反序列化客户端请求的报文
		req, err := DecodePacket(buf[:n])
		if err != nil {
			log.Println(err.Error())
			continue
		}

		//log.Println("recv request form ", addr.String(), req)

		req.T2 = T2
		req.T3 = TimeToTimeStamp(time.Now())

		// 将本地结构序列化
		newbuf, err := CodePacket(req)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 将序列化后的报文发送到客户端
		_, err = s.conn.WriteToUDP(newbuf, addr)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 请求成功次数统计
		s.request++
	}
}

// 申请时间服务器对象，并且初始化地址+端口
func NewNTPS(addr string) *NTPS {
	return &NTPS{addr: addr}
}

// 启动时间服务
func (s *NTPS) Start() error {

	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err2 := net.ListenUDP("udp", addr)
	if err2 != nil {
		return err2
	}

	s.wait.Add(1)
	s.conn = conn

	go msgProc(s)
	return nil
}

// 停止时间服务，并且释放资源
func (s *NTPS) Stop() {
	s.conn.Close()
	s.wait.Wait()
}
