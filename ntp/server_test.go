package ntp

import (
	"testing"
	"time"
)

func TestServerStart01(t *testing.T) {
	for i := 0; i < 2; i++ {
		ntps := NewNTPS("localhost", "1234")
		err := ntps.Start()
		if err != nil {
			t.Error(err.Error())
			return
		}
		ntps.Stop()
	}
}

func TestTimeStamp01(t *testing.T) {
	tm1 := TimeStamp{Sec: 10000, Nsec: 1}
	tm2 := TimeStamp{Sec: 0, Nsec: 999999999}

	tm1.Add(tm2)

	if tm1.Sec != 10001 || tm1.Nsec != 0 {
		t.Error(tm1, tm2)
	}
}

func TestTimeStamp02(t *testing.T) {
	tm1 := TimeStamp{Sec: 10000, Nsec: 0}
	tm2 := TimeStamp{Sec: 0, Nsec: 999999999}

	tm1.Sub(tm2)

	if tm1.Sec != 9999 || tm1.Nsec != 1 {
		t.Error(tm1, tm2)
	}
}

func TestTimeStamp03(t *testing.T) {
	tm1 := TimeStamp{Sec: 10000, Nsec: 0}
	tm2 := TimeStamp{Sec: 99, Nsec: 999999999}

	tm1.Div(10)
	tm2.Div(10)

	if tm1.Sec != 1000 || tm1.Nsec != 0 {
		t.Error(tm1, tm2)
	}

	if tm2.Sec != 9 || tm2.Nsec != 999999999 {
		t.Error(tm1, tm2)
	}
}

func TestTimeToTimeStamp01(t *testing.T) {
	tm1 := time.Now()
	tmstamp := TimeToTimeStamp(tm1)

	if tmstamp.Sec != tm1.Unix() {
		t.Error(tm1, tmstamp)
	}

	if tmstamp.Nsec != int64(tm1.Nanosecond()) {
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
	tm1 := time.Now()
	tm2 := TimeStampToTime(offset, tm1)

	if tm2 != tm1 {
		t.Error(tm1, tm2)
	}
}

func TestTimeStampToTime03(t *testing.T) {
	var offset TimeStamp

	offset.Sec = -1
	offset.Nsec = 0

	tm1 := time.Now()
	tm2 := TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix()-1 || tm2.Nanosecond() != tm1.Nanosecond() {
		t.Error(tm1, tm2)
	}

	offset.Sec = 1
	offset.Nsec = 0

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix()+1 || tm2.Nanosecond() != tm1.Nanosecond() {
		t.Error(tm1, tm2)
	}

	offset.Sec = 0
	offset.Nsec = 1

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix() || tm2.Nanosecond() != tm1.Nanosecond()+1 {
		t.Error(tm1, tm2)
	}

	offset.Sec = 0
	offset.Nsec = -1

	tm2 = TimeStampToTime(offset, tm1)

	if tm2.Unix() != tm1.Unix() || tm2.Nanosecond() != tm1.Nanosecond()-1 {
		t.Error(tm1, tm2)
	}
}
