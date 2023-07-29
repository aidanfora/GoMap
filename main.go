package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aidan/gomap/workers"
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
	-s  syn|con		Indicate the scanning mode. 
				SYN scans on localhost only, and it might crash your PC if you scan too many ports.
	-ip <IP_Address>	Indicate the IP Address to be scanned
	-p  <Port Numbers>	Indicate the ports to be scanned 
				Can be specified as a range or as individual ports
	-w  <Worker Numbers>	Indicate the number of worker functions to be launched as goroutines
				An increase in number will result in decreased reliability of scans

	Example: Basic TCP Con Scan of the first 1024 ports on your localhost
	.\gomap.exe -s con -ip 127.0.0.1 -p 1-1024
	`)
}

// Main parses command line flags and determines what mode, IP and port numbers are to be scanned
func main() {
	start := time.Now()

	fmt.Print(banner)

	// Variables to hold command line flags
	var openPorts []int
	var portRange []int
	var mode string
	var ipAddr string
	var ports string
	var workNum int64
	var showHelp bool

	// Assign command line flags tp variables. Defaults to TCP scan if unspecified, with 5000 worker functions running as goroutines
	flag.StringVar(&mode, "s", "con", "Scanning mode (syn|con)")
	flag.StringVar(&ipAddr, "ip", "", "IP Address to scan")
	flag.StringVar(&ports, "p", "", "Ports to scan")
	flag.Int64Var(&workNum, "w", 5000, "Number of worker functions to run")
	flag.BoolVar(&showHelp, "h", false, "Show available flags")
	flag.Parse()

	if showHelp {
		printMenu()
		os.Exit(0)
	}

	// Parse and format port ranges received
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
					return
				}

				end, err := strconv.Atoi(rangeParts[1])
				if err != nil {
					fmt.Println("Invalid ending port")
					return
				}

				for i := start; i <= end; i++ {
					portRange = append(portRange, i)
				}
			} else {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					fmt.Println("Invalid port number given")
					return
				}
				portRange = append(portRange, port)
			}
		}
	}
	fmt.Printf("Scan Mode: %s | IP Address: %s | Total Ports: %s | Workers Running: %s\n\n", mode, ipAddr, fmt.Sprint(len(portRange)), fmt.Sprint(workNum))

	// Create two channels, one to send ports to the worker function, and one to receive results from it
	portsChan := make(chan int, workNum)
	resultsChan := make(chan int)

	// Initialise worker functions
	if mode == "syn" {
		for i := 0; i <= cap(portsChan); i++ {
			go workers.WorkerSyn(portsChan, resultsChan, ipAddr)
		}
	} else if mode == "con" {
		for i := 0; i <= cap(portsChan); i++ {
			go workers.WorkerCon(portsChan, resultsChan, ipAddr)
		}
	} else {
		fmt.Println("Wrong scan mode specified")
		return
	}

	// Send ports to be scanned to the worker functions via channel portChan
	go func() {
		for _, port := range portRange {
			portsChan <- port
		}
	}()

	// Wait for results to be received from the resultsChan
	for i := 1; i <= len(portRange); i++ {
		port := <-resultsChan
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(portsChan)
	close(resultsChan)
	sort.Ints(openPorts)

	fmt.Printf("Format:\n%-5s | Service/Protocol\n\n", "Port")
	for _, port := range openPorts {
		serviceName, exists := detailedlist[port]
		if exists {
			fmt.Printf("%-5d | Service/Protocol: %s\n", port, serviceName)
		} else {
			fmt.Printf("%-5d | Service/Protocol: Unknown\n", port)
		}
	}

	end := time.Now()
	elapsed := end.Sub(start)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("\nScan Runtime: %d min %d sec\n", minutes, seconds)
}
