# GoMap

This is a repository for my mini project on developing a simple port scanner in Go

- Built in Go Version 1.20.6
- Requires the [gopacket libary](https://github.com/google/gopacket) and [Npcap](https://npcap.com/) to perform TCP Syn Scans

Credits to [JustinTimperio's detailed list](https://github.com/JustinTimperio/gomap/blob/master/gomap_ports.go) of ports and their respective service/protocol names

Tested on Windows 11 OS

## Usage

The `-h` option can be used for more information
```
.\gomap.exe -h

        ██████╗  ██████╗ ███╗   ███╗ █████╗ ██████╗
        ██╔════╝ ██╔═══██╗████╗ ████║██╔══██╗██╔══██╗
        ██║  ███╗██║   ██║██╔████╔██║███████║██████╔╝
        ██║   ██║██║   ██║██║╚██╔╝██║██╔══██║██╔═══╝
        ╚██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║
        ╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝

    A mini project to learn more about concurrency in Go


        Flags:
        -s  syn|con             Indicate the scanning mode.
                                SYN scans on localhost only, and it might crash your PC if you scan too many ports.
        -ip <IP_Address>        Indicate the IP Address to be scanned
        -p  <Port Numbers>      Indicate the ports to be scanned
                                Can be specified as a range or as individual ports
        -w  <Worker Numbers>    Indicate the number of worker functions to be launched as goroutines
                                An increase in number will result in decreased reliability of scans

        Example: Basic TCP Con Scan of the first 1024 ports on your localhost
        .\gomap.exe -s con -ip 127.0.0.1 -p 1-1024
```