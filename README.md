# GoMap

This is a repository for my mini project on developing a simple port scanner in Go

- Built in Go Version 1.20.6

## Usage

The `help` command can be used for more information
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
        -s  tcp|udp             Indicate the scanning mode
        -ip <IP_Address>        Indicate the IP Address to be scanned
        -p  <Port Numbers>      Indicate the ports to be scanned 
                                Can be specified as a range or as individual ports
        -w  <Worker Numbers>    Indicate the number of worker functions to be launched as goroutines
                                An increase in number will result in decreased reliability of scans

        Example: Basic TCP Scan of the first 1024 ports on your localhost
        .\gomap.exe -s tcp -ip 127.0.0.1 -p 1-1024
```