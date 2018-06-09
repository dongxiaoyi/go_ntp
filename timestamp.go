package ntp

import (
	"time"
)

type TimeStamp struct {
	NanoSecond int64 // 纳秒
}

// TimeStamp 时间结构转换 Time 类型
func TimeStampToTime(offset TimeStamp, now time.Time) time.Time {
	return now.Add(time.Duration(offset.NanoSecond))
}

// 将本地时间转换为 TimeStamp 结构
func TimeToTimeStamp(now time.Time) (tm TimeStamp) {
	tm.NanoSecond = int64(time.Second)*now.Unix() + int64(now.Nanosecond())
	return
}

// 获取绝对时间
func (t *TimeStamp) Abs() int64 {
	if t.NanoSecond < 0 {
		return -t.NanoSecond
	}
	return t.NanoSecond
}

// 获取绝对值（正负1）
func (t *TimeStamp) AbsValue() int64 {
	if t.NanoSecond < 0 {
		return -1
	}
	return 1
}

// TimeStamp时间sub操作
func (t *TimeStamp) Sub(s TimeStamp) TimeStamp {
	t.NanoSecond = t.NanoSecond - s.NanoSecond
	return *t
}

// TimeStamp时间add操作
func (t *TimeStamp) Add(a TimeStamp) TimeStamp {
	t.NanoSecond = t.NanoSecond + a.NanoSecond
	return *t
}

// TimeStamp时间除操作
func (t *TimeStamp) Div(d int) TimeStamp {
	t.NanoSecond = t.NanoSecond / int64(d)
	return *t
}

// 求TimeStamp时间平均值
func TimeStampAverage(d []TimeStamp) TimeStamp {
	var result TimeStamp
	if len(d) == 0 {
		return result
	}
	var sum TimeStamp
	for _, v := range d {
		sum.Add(v)
	}
	return sum.Div(len(d))
}
