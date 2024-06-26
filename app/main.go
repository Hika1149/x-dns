package main

import (
	"encoding/hex"
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

		hexStr := hex.EncodeToString(buf[:size])
		fmt.Printf("Hex data: %s\n", hexStr)

		//Convert the received data to binary
		binaryData := make([]int, len(buf[:size])*8)
		for i, b := range buf[:size] {
			for j := 0; j < 8; j++ {
				binaryData[i*8+j] = int((b >> (7 - j)) & 1)
			}
		}
		//Set the QR bit to 1
		binaryData[16] = 1

		response := bitsToByte(binaryData)

		//Convert the response to hex just for checking
		responseInHex := hex.EncodeToString(response)
		fmt.Println("Response in hex: ", responseInHex)

		//
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
