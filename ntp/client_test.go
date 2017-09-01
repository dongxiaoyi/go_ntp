package ntp

import (
	"testing"
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

		if result.NetDelay.Sec > 0 {
			t.Error(result)
		}

		if result.Offset.Sec > 0 {
			t.Error(result)
		}
	}
}

func TestClientEnd(t *testing.T) {
	ntps.Stop()
}
