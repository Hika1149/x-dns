package main

import (
	"fmt"
	"net"
	// Uncomment this block to pass the first stage
	// "net"
)

func bitsToByte(bits []int) []byte {
	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bits); i += 8 {
		for j := 0; j < 8; j++ {
			bytes[i/8] |= byte(bits[i+j]) << (7 - j)
		}
	}
	return bytes
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		response := make([]byte, size)

		copy(response, buf[:size])
		//binary.BigEndian.PutUint16(response[0:2], uint16(1234))
		response[2] = flipIndicator(response[2]) // set qr bit to 1

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
