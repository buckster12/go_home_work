package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Printf("current time: %s\n", time.Now().UTC())
	fmt.Printf("exact time: %s\n", ntpTime.UTC())
}
