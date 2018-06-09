package ntp

import (
	"testing"
	"time"
)

var ntps *NTPS

func TestClientBegin(t *testing.T) {
	ntps = NewNTPS("", "12345")
	err := ntps.Start()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestClient01(t *testing.T) {
	ntpc := NewNTPC("localhost", "12345")

	for i := 0; i < 10; i++ {

		result, err := ntpc.Sync(10)
		if err != nil {
			t.Error(err.Error())
			break
		}

		if result.NetDelay.NanoSecond > int64(time.Second) {
			t.Error(result)
		}

		if result.Offset.NanoSecond > int64(time.Second) {
			t.Error(result)
		}
	}
}

func TestClientEnd(t *testing.T) {
	ntps.Stop()
}
