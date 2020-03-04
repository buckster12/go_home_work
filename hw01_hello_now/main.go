package main

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime := time.Now().UTC()
	ntpTime, _ := ntp.Time("0.beevik-ntp.pool.ntp.org")
	dateFormat := "2006-01-02 15:04:05 +0000 UTC"
	fmt.Printf("current time: %s\n", currentTime.Format(dateFormat))
	fmt.Printf("exact time: %s\n", ntpTime.UTC().Format(dateFormat))
}
