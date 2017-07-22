package main

import (
	"fmt"
	"go_ntp/ntp"
)

func main() {
	ntpc := ntp.NewNTPC("time.nist.gov", "123")

	ntpc.Config(1, 30)

	for i := 0; i < 10000; i++ {
		err := ntpc.Sync()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	return
}
