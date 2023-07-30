package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
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
	-ip <IP Address>	Indicate the IP Address to be scanned
	-p  <Port Numbers>	Indicate the ports to be scanned 
				Can be specified as a range or as individual ports
	-w  <Worker Numbers>	Indicate the number of worker functions to be launched as goroutines
				An increase in number above 15000 may result in decreased reliability of scans

	Example: Basic TCP Connect Scan of the first 1024 ports on your localhost
	.\gomap.exe -ip 127.0.0.1 -p 1-1024
	`)
}

// main parses the command line flags and determines the IP & ports to be scanned, and number of worker functions to run
func main() {
	start := time.Now()

	fmt.Print(banner)

	// Variables to hold command line flags
	var openPorts []int
	var portRange []int
	var ipAddr string
	var ports string
	var workNum int64
	var showHelp bool

	// Assign command line flags to variables. Defaults to 5000 worker functions running as goroutines
	flag.StringVar(&ipAddr, "ip", "", "IP Address to scan")
	flag.StringVar(&ports, "p", "", "Ports to scan")
	flag.Int64Var(&workNum, "w", 5000, "Number of worker functions to run")
	flag.BoolVar(&showHelp, "h", false, "Show available flags")
	flag.Parse()

	if showHelp {
		printMenu()
		os.Exit(0)
	}

	portFormat(ports, &portRange)

	// Prints some details about the scan, as well as the possible OS of the host
	fmt.Printf("Scan Mode: %s | IP Address: %s | Total Ports: %s | Workers Running: %s\n", "tcp", ipAddr, fmt.Sprint(len(portRange)), fmt.Sprint(workNum))
	detected, err := osDetection(ipAddr)
	if err != nil {
		fmt.Println("Error detecting OS", err.Error())
		return
	}
	fmt.Printf("Possible OS Detected: %s\n\n", detected)

	// Create two channels, one to send ports to the worker function, and one to receive results from it
	portsChan := make(chan int, workNum)
	resultsChan := make(chan int)

	// Initialise worker functions
	for i := 0; i <= cap(portsChan); i++ {
		go worker(portsChan, resultsChan, ipAddr)
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

	// Close channels and sort the list of open ports
	close(portsChan)
	close(resultsChan)
	sort.Ints(openPorts)

	// Print output
	fmt.Printf("Format:\n%-5s | Service/Protocol\n\n", "Port")
	for _, port := range openPorts {
		serviceName, exists := detailedlist[port]
		if exists {
			fmt.Printf("%-5d | %s\n", port, serviceName)
		} else {
			fmt.Printf("%-5d | Unknown\n", port)
		}
	}

	// Calculate the scan runtime
	end := time.Now()
	elapsed := end.Sub(start)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("\nScan Runtime: %d min %d sec\n", minutes, seconds)
}
