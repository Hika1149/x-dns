package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"github.com/codecrafters-io/dns-server-starter-go/app/udp"
	"net"
	// Uncomment this block to pass the first stage
	// "net"
)

func buildResponsePacket(request *dns.DNSPacket) *dns.DNSPacket {

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

	// set each answer
	for _, q := range packet.Questions {
		packet.AddAnswer(&dns.Record{
			Name:   q.Name,
			Type:   1,
			Class:  1,
			TTL:    60,
			Length: 4,
			Data:   "8.8.8.8",
		})
	}
	return packet

}

func buildResponsePacketByForwardingServer(request *dns.DNSPacket, resolverIp string) *dns.DNSPacket {
	// Assumptions
	//- When you receive multiple questions in the question section you will need to split it into two DNS packets
	//- Send them to this resolver then merge the response in a single packet.

	// init res
	res := dns.NewDNSPacket()
	res.Header = request.Header
	res.Header.Response = true
	res.Questions = request.Questions

	// split request
	reqs := make([][]byte, 0)
	for _, q := range request.Questions {
		packet := dns.NewDNSPacket()
		packet.Header = request.Header
		packet.Header.QuestionCount = 1
		packet.Questions = []*dns.DNSQuestion{q}

		reqByte, err := packet.ToByte()
		if err != nil {
			fmt.Println("Failed to convert request to bytes:", err)
			continue
		}
		reqs = append(reqs, reqByte)
	}

	// send each req
	for _, reqByte := range reqs {
		// send to resolver server
		result, err := udp.Dial(resolverIp, reqByte)
		if err != nil {
			fmt.Println("Failed to send request:", err)
			return nil
		}

		resPacket := dns.NewDNSPacket().FromByte(result)

		// append to res
		res.Header.ResCode = resPacket.Header.ResCode
		for _, answer := range resPacket.Answers {
			res.AddAnswer(answer)
		}
	}

	return res

}

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

	resolverIp := flag.String("resolver", "", "Forwarding server IP")
	flag.Parse()

	//
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		fmt.Println("buf size: ", size)
		fmt.Println("resolver ip: ", *resolverIp)

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
		var packet *dns.DNSPacket

		if *resolverIp != "" {
			packet = buildResponsePacketByForwardingServer(request, *resolverIp)
		} else {
			packet = buildResponsePacket(request)
		}

		// write to response buffer

		resByte, err := packet.ToByte()
		if err != nil {
			fmt.Println("Failed to convert response to bytes:", err)
			break
		}

		//
		fmt.Println("resByte: ", resByte, len(resByte))

		_, err = udpConn.WriteToUDP(resByte, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func flipIndicator(b byte) byte {
	return b | 0b10000000
}
