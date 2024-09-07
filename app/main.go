package main

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
	// Uncomment this block to pass the first stage
	// "net"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	//
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	//
	buf := make([]byte, 512)
	//
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		reqBuffer := buffer.BytePacketBuffer{
			Buffer: buf[:size],
			Pos:    0,
		}

		var request = dns.NewDNSPacket()

		err = request.FromBuffer(reqBuffer)

		if err != nil {
			fmt.Printf("FromBuffer failed: %v\n", err)
			break
		}

		// response packet
		packet := dns.NewDNSPacket()

		// init from request
		packet.Header = request.Header
		packet.Questions = request.Questions

		// set response header
		packet.Header.Response = true
		packet.Header.AuthoritativeAnswer = false
		packet.Header.TruncatedMsg = false
		packet.Header.RecursionAvailable = false
		packet.Header.Z = 0
		if request.Header.OpCode == 0 {
			packet.Header.ResCode = 0
		} else {
			packet.Header.ResCode = 4
		}

		packet.AddAnswer(&dns.Record{
			Name:   "codecrafters.io",
			Type:   1,
			Class:  1,
			TTL:    60,
			Length: 4,
			Data:   "8.8.8.8",
		})

		// write to response buffer
		resBuffer := buffer.NewBytePacketBuffer()
		err = packet.Write(resBuffer)
		if err != nil {
			fmt.Println("Failed to write response:", err)
			break
		}
		fmt.Println("debug: ", resBuffer.Buffer[:resBuffer.Pos], resBuffer.Pos)

		//
		_, err = udpConn.WriteToUDP(resBuffer.Buffer[:resBuffer.Pos], source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func flipIndicator(b byte) byte {
	return b | 0b10000000
}
