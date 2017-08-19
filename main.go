package main

import (
	"fmt"
	"go_ntp/ntp"
	"os"
	"time"
)

func Server() {

	ntps := ntp.NewNTPS("", "1234")

	ntps.Start()

	for {
		time.Sleep(1 * time.Second)
	}

	ntps.Stop()
}

func Client() {
	ntpc := ntp.NewNTPC("localhost", "1234")

	ntpc.Config(2, 30)

	for i := 0; i < 10000; i++ {
		err := ntpc.Sync()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: <-s/-c>")
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
