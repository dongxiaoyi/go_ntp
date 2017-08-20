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
	DEFAULT_PACKET_SIZE = 10 * 8
)

type TimeStamp struct {
	Sec  int64
	Nsec int64
}

type Packet struct {
	Version   uint64
	RequestId uint64
	T1        TimeStamp // T1 客户端发送请求的时间
	T2        TimeStamp // T2 服务器接收请求的时间
	T3        TimeStamp // T3 服务器答复时间
	T4        TimeStamp // T4 客户端接收答复时间
}

type NTPS struct {
	request int64
	addr    string
	conn    *net.UDPConn
	wait    sync.WaitGroup
}

func GetTimeStamp() TimeStamp {
	var tm TimeStamp
	now := time.Now()

	tm.Sec = now.Unix()
	tm.Nsec = int64(now.Nanosecond())

	return tm
}

func (t *TimeStamp) Sub(s TimeStamp) TimeStamp {
	t.Sec = t.Sec - s.Sec

	if t.Nsec < s.Nsec {
		t.Nsec = int64(time.Second) + t.Nsec - s.Nsec
		t.Sec--
	} else {
		t.Nsec = t.Nsec - s.Nsec
	}

	return *t
}

func (t *TimeStamp) Add(a TimeStamp) TimeStamp {
	t.Sec = t.Sec + a.Sec
	t.Nsec = t.Nsec + a.Nsec

	if t.Nsec > int64(time.Second) {
		t.Nsec = t.Nsec - int64(time.Second)
		t.Sec++
	}

	return *t
}

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

func DecodePacket(buf []byte) (rsp Packet, err error) {

	iobuf := bytes.NewReader(buf)
	err = binary.Read(iobuf, binary.BigEndian, &rsp)

	//log.Println("RSP: ", rsp)
	//log.Println("RECV_BUF:", len(buf), buf)

	return
}

func msgProc(s *NTPS) {

	defer s.wait.Done()

	var buf [4096]byte

	for {
		n, addr, err := s.conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println("socket disable. ", s.conn.RemoteAddr())
			return
		}

		T2 := GetTimeStamp()

		if n != DEFAULT_PACKET_SIZE {
			log.Println("recv a packet not recognized ", len(buf), buf[0:n])
			continue
		}

		req, err := DecodePacket(buf[:n])
		if err != nil {
			log.Println(err.Error())
			continue
		}

		//log.Println("recv request form ", addr.String(), req)

		req.T2 = T2
		req.T3 = GetTimeStamp()

		newbuf, err := CodePacket(req)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		_, err = s.conn.WriteToUDP(newbuf, addr)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

func NewNTPS(ip, port string) *NTPS {
	return &NTPS{addr: ip + ":" + port}
}

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

func (s *NTPS) Stop() {
	s.conn.Close()
	s.wait.Wait()
}
