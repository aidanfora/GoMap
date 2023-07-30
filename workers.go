package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// worker performs scanning using TCP Connect Scans which establish a full connection with the target
// It receive ports via portChan, performs the scan and collects the result, and sends result into resultsChan
func worker(portsChan <-chan int, resultsChan chan<- int, ipAddr string) {
	for p := range portsChan {
		address := fmt.Sprintf("%s:%d", ipAddr, p)
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)

		if err != nil {
			resultsChan <- 0
			continue
		} else {
			conn.Close()
			resultsChan <- p
		}
	}
}

// portFormat performs formatting of the ports and port ranges passed as arguments
func portFormat(ports string, portRange *[]int) {
	if ports == "-" {
		for i := 1; i <= 65535; i++ {
			*portRange = append(*portRange, i)
		}
	} else {
		portStrings := strings.Split(ports, ",")
		for _, portStr := range portStrings {
			if strings.Contains(portStr, "-") {
				rangeParts := strings.Split(portStr, "-")
				start, err := strconv.Atoi(rangeParts[0])
				if err != nil {
					fmt.Println("Invalid starting port")
					return
				}

				end, err := strconv.Atoi(rangeParts[1])
				if err != nil {
					fmt.Println("Invalid ending port")
					return
				}

				for i := start; i <= end; i++ {
					*portRange = append(*portRange, i)
				}
			} else {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					fmt.Println("Invalid port number given")
					return
				}
				*portRange = append(*portRange, port)
			}
		}
	}

}
