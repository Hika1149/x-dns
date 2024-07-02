package main

import (
	"bytes"
	"encoding/binary"
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

		//var msg DNSMessage
		//
		//err = binary.Read(bytes.NewReader(buf[:size]), binary.BigEndian, &msg)

		var msg DNSMessage
		err = binary.Read(bytes.NewReader(buf[:size]), binary.BigEndian, &msg)
		fmt.Printf("read binary failed: %v\n", err)

		fmt.Printf("msg: %v %v ", msg.Header, len(msg.Question))

		response := msg.ToBytes()

		//
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func flipIndicator(b byte) byte {
	return b | 0b10000000
}
