package ntp

import (
	"testing"
	"time"
)

func TestSetTimeToOs01(t *testing.T) {
	tm := time.Now()
	err := SetTimeToOs(tm)
	if err != nil {
		t.Error("set time failed!" + err.Error())
		return
	}
}

func TestSetTimeToOs02(t *testing.T) {

	tmOld := time.Now()
	tmNew := tmOld

	tmNew = tmNew.Add(time.Hour)

	err := SetTimeToOs(tmNew)
	if err != nil {
		t.Error("set time failed!" + err.Error())
		return
	}

	time.Sleep(1000 * time.Millisecond)

	time.Now()

	tmNew2 := time.Now()

	if tmNew2.Sub(tmOld) < time.Second {
		t.Error(tmNew2, tmOld, tmNew)
	}

	t.Log(tmNew2, tmOld, tmNew)

	err = SetTimeToOs(tmOld)
	if err != nil {
		t.Error("set time failed!" + err.Error())
		return
	}
}
