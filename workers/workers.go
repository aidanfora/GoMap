package workers

import (
	"fmt"
	"net"
	"time"
)

// Worker functions receive ports via a channel and prints to console whether port is opened or closed

// WorkerCon performs scanning using TCP Connect Scans which establish a full connection with the target
func Worker(portsChan <-chan int, resultsChan chan<- int, ipAddr string) {
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
