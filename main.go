package main

import (
	"fmt"
	"go_ntp/ntp"
)

func main() {
	ntpc := ntp.NewNTPC("sim.ntp.org.cn", "123")

	err := ntpc.Sync()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}
