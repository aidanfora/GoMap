package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const banner string = `
	██████╗  ██████╗ ███╗   ███╗ █████╗ ██████╗ 
	██╔════╝ ██╔═══██╗████╗ ████║██╔══██╗██╔══██╗
	██║  ███╗██║   ██║██╔████╔██║███████║██████╔╝
	██║   ██║██║   ██║██║╚██╔╝██║██╔══██║██╔═══╝ 
	╚██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║     
	╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     

    A mini project to learn more about concurrency in Go

`

// printMenu displays a help menu showing all usable flags
func printMenu() {
	fmt.Println(`
	Flags:
	-s  tcp|udp		Indicate the scanning mode
	-ip <IP_Address>	Indicate the IP Address to be scanned
	-p  <Port Numbers>	Indicate the ports to be scanned 
				Can be specified as a range or as individual ports
	-w  <Worker Numbers>	Indicate the number of worker functions to be launched as goroutines
				An increase in number will result in decreased reliability of scans

	Example: Basic TCP Scan of the first 1024 ports on your localhost
	.\gomap.exe -s tcp -ip 127.0.0.1 -p 1-1024
	`)
}

// Worker receives ports via a channel and prints to console whether port is opened or closed
func worker(portsChan <-chan int, resultsChan chan<- int, mode string, ipAddr string) {
	for p := range portsChan {
		address := ipAddr + fmt.Sprintf(":%d", p)
		conn, err := net.DialTimeout(mode, address, 5*time.Second)
		if err != nil {
			resultsChan <- 0
			continue
		} else {
			conn.Close()
			resultsChan <- p
		}
	}
}

// Main parses command line flags and determines what mode, IP and port numbers are to be scanned
func main() {
	fmt.Print(banner)

	// Variables to hold command line flags
	var openPorts []int
	var portRange []int
	var mode string
	var ipAddr string
	var ports string
	var showHelp bool
	var workers int64

	// Set up command line flags. Defaults to TCP scan if unspecified
	flag.StringVar(&mode, "s", "tcp", "Scan mode (udp|tcp)")
	flag.StringVar(&ipAddr, "ip", "", "IP Address to scan")
	flag.StringVar(&ports, "p", "", "Ports to scan")
	flag.Int64Var(&workers, "w", 1000, "Number of worker functions to run")
	flag.BoolVar(&showHelp, "h", false, "Show available flags")
	flag.Parse()

	if showHelp {
		printMenu()
		os.Exit(0)
	}

	// Parse and format ports received
	if ports == "-" {
		for i := 1; i <= 65535; i++ {
			portRange = append(portRange, i)
		}
	} else {
		portStrings := strings.Split(ports, ",")
		for _, portStr := range portStrings {
			if strings.Contains(portStr, "-") {
				rangeParts := strings.Split(portStr, "-")
				start, err := strconv.Atoi(rangeParts[0])
				if err != nil {
					fmt.Println("Invalid starting port")
				}

				end, err := strconv.Atoi(rangeParts[1])
				if err != nil {
					fmt.Println("Invalid ending port")
				}

				for i := start; i <= end; i++ {
					portRange = append(portRange, i)
				}
			} else {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					fmt.Println("Invalid port number given")
				}
				portRange = append(portRange, port)
			}
		}
	}
	fmt.Printf("Scan Mode: %s | IP Address: %s | Total Ports: %s | Workers Running: %s\n\n", mode, ipAddr, fmt.Sprint(len(portRange)), fmt.Sprint(workers))

	// Create two channels, one to send ports to the worker function, and one to receive results from it
	portsChan := make(chan int, workers)
	resultsChan := make(chan int)

	for i := 0; i <= cap(portsChan); i++ {
		go worker(portsChan, resultsChan, mode, ipAddr)
	}

	go func() {
		for _, port := range portRange {
			portsChan <- port
		}
	}()

	start := time.Now()

	for i := 1; i <= len(portRange); i++ {
		port := <-resultsChan
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(portsChan)
	close(resultsChan)
	sort.Ints(openPorts)

	for _, port := range openPorts {
		fmt.Printf("%d is open\n", port)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("\nScan Runtime: %d min %d sec\n", minutes, seconds)
}
