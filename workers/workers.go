package workers

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func halfOpenScan(ipAddr string, port int) (bool, error) {
	handle, err := pcap.OpenLive("\\Device\\NPF_Loopback", 1600, false, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Encountered error at pcap.OpenLive: %s", err.Error())
		return false, err
	}

	// Construct the TCP SYN packet
	options := gopacket.SerializeOptions{
		ComputeChecksums: true, // Compute checksum for IP, TCP and UDP layers.
		FixLengths:       true, // Automatically computes lengths.
	}
	buffer := gopacket.NewSerializeBuffer()

	// Assign a random source port between 1025 and 65535
	srcPort := rand.Intn(65535-1025) + 1025
	srcIP := net.IP{0, 0, 0, 0}

	gopacket.SerializeLayers(buffer, options,
		&layers.Ethernet{},
		&layers.IPv4{
			Version:  4,
			TTL:      64,
			SrcIP:    srcIP,
			DstIP:    net.ParseIP(ipAddr),
			Protocol: layers.IPProtocolTCP,
		},
		&layers.TCP{
			SrcPort: layers.TCPPort(srcPort),
			DstPort: layers.TCPPort(port),
			SYN:     true,
			Seq:     11050, // Random sequence number
			Window:  14600, // Arbitrary chosen
		},
		gopacket.Payload([]byte{}),
	)
	outgoingPacket := buffer.Bytes()

	// Send the SYN packet
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return false, err
	}

	// Wait for the SYN-ACK or RST-ACK response
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Parse the packet
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if tcpLayer != nil && ipLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			ip, _ := ipLayer.(*layers.IPv4)

			// Check that the packet is a response to our SYN
			if tcp.DstPort == layers.TCPPort(srcPort) && ip.DstIP.Equal(srcIP) {
				// If it's a SYN-ACK packet, return true
				if tcp.SYN && tcp.ACK {
					return true, nil
				}
				// If it's a RST-ACK packet, return false
				if tcp.RST && tcp.ACK {
					return false, nil
				}
			}
		}
	}

	// If we got no response at all, consider the port filtered
	return false, nil
}

// Worker functions receive ports via a channel and prints to console whether port is opened or closed

// WorkerSyn performs scanning using TCP Syn Scans which establish a half connection with the target
func WorkerSyn(portsChan <-chan int, resultsChan chan<- int, ipAddr string) {
	for p := range portsChan {
		isOpen, err := halfOpenScan(ipAddr, p)
		if err != nil || !isOpen {
			resultsChan <- 0
			continue
		} else {
			resultsChan <- p
		}
	}
}

// WorkerCon performs scanning using TCP Connect Scans which establish a full connection with the target
func WorkerCon(portsChan <-chan int, resultsChan chan<- int, ipAddr string) {
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
