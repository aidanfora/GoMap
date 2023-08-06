# GoMap

This is a repository for my mini project on developing a simple port scanner in Go

- Built in Go Version 1.20.6

- Tested on Windows 11 OS

- Credits to [JustinTimperio's detailed list](https://github.com/JustinTimperio/gomap/blob/master/gomap_ports.go) of ports and their respective service/protocol names

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
        -ip <IP Address>        Indicate the IP Address to be scanned
        -p  <Port Numbers>      Indicate the ports to be scanned
                                Can be specified as a range or as individual ports
        -w  <Worker Numbers>    Indicate the number of worker functions to be launched as goroutines
                                An increase in number above 15000 may result in decreased reliability of scans   

        Example: Basic TCP Connect Scan of the first 1024 ports on your localhost
        .\gomap.exe -ip 127.0.0.1 -p 1-1024
```