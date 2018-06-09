package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lixiangyun/go_ntp"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: < IP:PORT >")
		return
	}

	// Listens to addresses and ports
	ntps := ntp.NewNTPS(args[1])

	// Start the service, then coroutine process is created in the background.
	ntps.Start()
	for {
		time.Sleep(60 * time.Second)
	}

	// Stop the service.
	ntps.Stop()
}
