package ntp

import (
	"testing"
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
