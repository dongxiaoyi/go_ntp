package ntp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	NTPC_VERSION            = 3   // NTPv3
	NTPC_LEAP_NOTINSYNC     = 0x3 //
	NTPC_MODE_CLIENT        = 3
	NTPC_STRATUM_PKT_UNSPEC = 0
	NTPC_MINPOLL            = 4
	NTPC_PRECISION          = 0xfa
	NTPC_FP_SECOND          = 0x00010000
	NTPC_SHIFT_BITS_SIXTEEN = 16
	NTPC_JAN_1970           = 0x83aa7e80 // 3600s*24h*(365d*70y+17d)

	NTPC_DISTANCE = NTPC_FP_SECOND
	NTPC_DISP     = NTPC_FP_SECOND

	NTPC_DEFAULT_PORT       = "123"
	NTPC_DEFAULT_TIMEOUT    = 1
	NTPC_DEFAULT_RETRYTIMES = 10

	NTPC_NTP_PACKET_SIZE = 48
	NTPC_NTP_AUTH_SIZE   = 160
)

type localTimeStamp struct {
	Sec  int32
	Farc int32
}

type ntpTimeStamp struct {
	Sec  uint32
	Farc uint32
}

/*
·  LI：跳跃指示器，警告在当月最后一天的最终时刻插入的迫近闺秒（闺秒）。
·  VN：版本号。
·  Mode：模式。该字段包括以下值：0－预留；1－对称行为；3－客户机；4－服务器；5－广播；6－NTP 控制信息
·  Stratum：对本地时钟级别的整体识别。
·  Poll：有符号整数表示连续信息间的最大间隔。
·  Precision：有符号整数表示本地时钟精确度。
·  Root Delay：有符号固定点序号表示主要参考源的总延迟，很短时间内的位15到16间的分段点。
·  Root Dispersion：无符号固定点序号表示相对于主要参考源的正常差错，很短时间内的位15到16间的分段点。
·  Reference Identifier：识别特殊参考源。
·  Originate Timestamp：这是向服务器请求分离客户机的时间，采用64位时标（Timestamp）格式。
·  Receive Timestamp：这是向服务器请求到达服务器的时间，采用64位时标（Timestamp）格式。
·  Transmit Timestamp：这是向客户机答复分离服务器的时间，采用64位时标（Timestamp）格式。
*/
type ntpPacket struct {
	LiVnMode            byte
	Stratum             byte
	Poll                byte
	Precision           byte
	RootDelay           uint32
	RootDispersion      uint32
	ReferenceIdentifier uint32
	ReferenceTimestamp  ntpTimeStamp
	OriginateTimestamp  ntpTimeStamp
	ReceiveTimestamp    ntpTimeStamp
	TransmitTimestamp   ntpTimeStamp
}

/*
Authenticator（Optional）：
当实现了 NTP 认证模式,主要标识符和信息数字域就包括已定义的信息认证代码（MAC）信息。
*/
type ntpAuthenticator struct {
	KeyIdentifier [32]byte
	MessageDigest [128]byte
}

type NTPC struct {
	IP   string // NTPS IP Addr
	PORT string // NTPS PORT

	TIMEOUT    int // wait resp timeout
	RETRYTIMES int // retry times
}

func makeLiVnMode(Li, Vn, Mode byte) byte {
	var temp byte

	temp = (Li << 6) & 0xC0
	temp |= (Vn << 3) & 0x38
	temp |= (Mode & 0x7)

	return temp
}

func timeToTimeStamp(t time.Time) (tmsp ntpTimeStamp) {
	var fraction float64

	fraction = float64(t.Nanosecond() / 1000)
	fraction = fraction * float64(2^32) / float64(10^6)

	tmsp.Sec = uint32(t.Unix() + NTPC_JAN_1970)
	tmsp.Farc = uint32(fraction)

	return
}

func timeStampToTime(tmsp ntpTimeStamp) time.Time {
	var nsec float64

	nsec = float64(tmsp.Farc) * float64(10^6) / float64(2^32)

	return time.Unix(int64(tmsp.Sec), int64(nsec))
}

func NewNTPC(ip, port string) *NTPC {

	var ntpc = NTPC{IP: ip, PORT: port}

	ntpc.TIMEOUT = NTPC_DEFAULT_TIMEOUT
	ntpc.RETRYTIMES = NTPC_DEFAULT_RETRYTIMES

	return &ntpc
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func makeNtpPacket(t time.Time) ntpPacket {
	var req ntpPacket

	req.LiVnMode = makeLiVnMode(NTPC_LEAP_NOTINSYNC, NTPC_VERSION, NTPC_MODE_CLIENT)
	req.Stratum = NTPC_STRATUM_PKT_UNSPEC
	req.Poll = NTPC_MINPOLL
	req.Precision = NTPC_PRECISION
	req.RootDelay = NTPC_DISTANCE >> NTPC_SHIFT_BITS_SIXTEEN
	req.RootDispersion = NTPC_DISP >> NTPC_SHIFT_BITS_SIXTEEN
	req.ReferenceIdentifier = 0
	req.ReferenceTimestamp = timeToTimeStamp(t)

	return req
}

func sendNtpPacket(s net.Conn, req ntpPacket) error {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, req)
	if err != nil {
		fmt.Println("Serialization failed.", err.Error())
		return err
	}

	fmt.Println("REQ: ", req)
	fmt.Println("SENDBUF: ", buf.Len(), buf.Bytes())

	_, err = s.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Send Pkt NTPS failed.", err.Error())
		return err
	}

	return nil
}

func recvNtpPacket(s net.Conn) (rsp ntpPacket, err error) {
	var buf [1024]byte
	var cnt int

	cnt, err = s.Read(buf[0:])
	if err != nil {
		return
	}

	if cnt < NTPC_NTP_PACKET_SIZE {
		err = errors.New("Recv Pkt From NTPS failed.")
		return
	}

	fmt.Println("RECV_BUF:", cnt, buf[:cnt])

	iobuf := bytes.NewReader(buf[:cnt])
	err = binary.Read(iobuf, binary.BigEndian, &rsp)
	return
}

func calcDiffTime(req, rsp ntpPacket) (dly, off ntpTimeStamp) {
	var t1, t2, t3, t4 ntpTimeStamp

	t1 = req.ReferenceTimestamp      // T1 客户端发送请求的时间
	t2 = rsp.ReceiveTimestamp        // T2 服务器接收请求的时间
	t3 = rsp.TransmitTimestamp       // T3 服务器答复时间
	t4 = timeToTimeStamp(time.Now()) // T4 客户端接收答复时间

	dly.Sec = (t2.Sec - t1.Sec) + (t4.Sec - t3.Sec) // 计算得出网络时延
	dly.Farc = (t2.Farc - t1.Farc) + (t4.Farc - t3.Farc)

	off.Sec = (t2.Sec - t1.Sec) + (t3.Sec - t4.Sec) // 计算本地与服务器时延
	off.Farc = (t2.Farc - t1.Farc) + (t3.Farc - t4.Farc)

	return
}

func (n *NTPC) Sync() error {
	var req, rsp ntpPacket
	var timeout time.Time

	socket, err := net.Dial("udp", n.IP+":"+n.PORT)
	if err != nil {
		return err
	}

	defer socket.Close()

	for i := 0; i < n.RETRYTIMES; i++ {

		timeout = time.Now()
		timeout = timeout.Add(time.Duration(n.TIMEOUT) * time.Second)

		err = socket.SetDeadline(timeout)
		if err != nil {
			return err
		}

		req = makeNtpPacket(time.Now())

		err = sendNtpPacket(socket, req)
		if err != nil {
			return err
		}

		rsp, err = recvNtpPacket(socket)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		break
	}

	fmt.Println("REQ:", req)
	fmt.Println("RSP:", rsp)

	dly, off := calcDiffTime(req, rsp)

	netdly := timeStampToTime(dly)
	offset := timeStampToTime(off)

	fmt.Println("NetDelay:", netdly.Nanosecond()/int(time.Microsecond))

	fmt.Println("Offser:", offset.Second())

	return nil
}
