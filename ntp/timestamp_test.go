package ntp

import (
	"testing"
	"time"
)

func TestTimeStamp01(t *testing.T) {
	tm1 := TimeStamp{NanoSecond: 1}
	tm2 := TimeStamp{NanoSecond: 999999999}

	tm1.Add(tm2)

	if tm1.NanoSecond != int64(time.Second) {
		t.Error(tm1, tm2)
	}
}

func TestTimeStamp02(t *testing.T) {
	tm1 := TimeStamp{NanoSecond: 1000000000}
	tm2 := TimeStamp{NanoSecond: 999999999}

	tm1.Sub(tm2)

	if tm1.NanoSecond != 1 {
		t.Error(tm1, tm2)
	}
}

func TestTimeStamp03(t *testing.T) {
	tm1 := TimeStamp{NanoSecond: 1000000000}
	tm2 := TimeStamp{NanoSecond: 999999999}

	tm1.Div(10)
	tm2.Div(10)

	if tm1.NanoSecond != 100000000 {
		t.Error(tm1, tm2)
	}

	if tm2.NanoSecond != 99999999 {
		t.Error(tm1, tm2)
	}
}

func TestTimeToTimeStamp01(t *testing.T) {
	tm1 := time.Now()
	tmstamp := TimeToTimeStamp(tm1)

	if tmstamp.NanoSecond != tm1.Unix()*int64(time.Second)+int64(tm1.Nanosecond()) {
		t.Error(tm1, tmstamp)
	}
}

func TestTimeStampToTime01(t *testing.T) {
	var offset TimeStamp
	tm1 := time.Now()
	tm2 := TimeStampToTime(offset, tm1)

	if tm2 != tm1 {
		t.Error(tm1, tm2)
	}
}

func TestTimeStampToTime02(t *testing.T) {
	var offset TimeStamp

	offset.NanoSecond = -int64(time.Second)

	tm1 := time.Now()
	tm2 := TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix()-1 || tm2.Nanosecond() != tm1.Nanosecond() {
		t.Error(tm1, tm2)
	}

	offset.NanoSecond = int64(time.Second)

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix()+1 || tm2.Nanosecond() != tm1.Nanosecond() {
		t.Error(tm1, tm2)
	}

	offset.NanoSecond = 1

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix() || tm2.Nanosecond() != tm1.Nanosecond()+1 {
		t.Error(tm1, tm2)
	}

	offset.NanoSecond = -1

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix() || tm2.Nanosecond() != tm1.Nanosecond()-1 {
		t.Error(tm1, tm2)
	}
}
