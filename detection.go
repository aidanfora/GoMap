package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

// ICMP structure for manually creating ICMP echo packets
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

// osDetection attempts to fingerprint the OS by checking the value of the TTL field from an ICMP Echo Reply
func osDetection(ipAddr string) (string, error) {
	conn, err := net.Dial("ip4:icmp", ipAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Create the ICMP Echo Request
	var icmp ICMP
	icmp.Type = 8
	icmp.Code = 0
	icmp.Identifier = 1
	icmp.SequenceNum = 1

	// Convert the ICMP struct into bytes, calculate the checksum, and rewrite it to a buffer
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return "", err
	}

	// Prepare to receive the ICMP Echo Reply
	recv := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second * 2)) // Unblock after 2 seconds of no reply
	if _, err := conn.Read(recv); err != nil {
		return "", err
	}

	// Checks if it is an IP Packet
	if recv[0] == 0x45 {
		ttl := recv[8]
		if ttl <= 64 {
			return "Linux", nil
		} else if ttl <= 128 {
			return "Windows", nil
		} else if ttl <= 255 {
			return "Unix-like/Cisco Router/IoT Device", nil
		}
	}

	return "Unknown", nil
}

// checkSum calculates ICMP packet checksum
func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	if length > 0 {
		sum += uint32(data[index]) << 8
	}

	sum += (sum >> 16)

	return uint16(^sum)
}
