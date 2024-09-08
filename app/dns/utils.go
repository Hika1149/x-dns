package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
	"strings"
)

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func encodeDomainName(domain string) []byte {
	// domain name encoded as a sequence of labels
	// where each label consists of a length octet followed by that
	// many octets
	// the domain name is terminated by a zero length octet
	// the domain name is case-insensitive

	res := make([]byte, 0)

	splitS := strings.Split(domain, ".")

	for _, s := range splitS {

		l := len(s)

		b := []byte(s)

		//fmt.Println("label length", byte(l), "label", b)
		res = append(res, byte(l))
		res = append(res, b...)
	}
	// terminating zero length octet
	res = append(res, 0x00)
	return res

}

func DecodeDomainName(buffer buffer.BufferReader) (string, error) {

	// domain name encoded as a sequence of labels
	// where each label consists of a length octet followed by that
	// number of octets.

	// <length><content>

	res := ""
	var posRestored uint16
	for {
		length, err := buffer.ReadU8()
		if err != nil {
			return "", err
		}

		if length == 0 {
			break
		}

		// check if this is a pointer
		//if (length & 0b11000000) == 1 {
		isPointer := (length & 0xC0) == 0xC0

		if isPointer {
			// The pointer takes the form of a two octet sequence:
			b, err := buffer.ReadU8()
			if err != nil {
				return "", err
			}

			offset := uint16(length & 0x3F)
			offset = offset << 8
			offset = offset | uint16(b)

			posRestored = buffer.Position()
			buffer.SetPosition(offset)
			fmt.Printf("record detect pointer nextPos=%v offset=%v bufferPos=%v\n", posRestored, offset, buffer.Position())
			continue

		}

		if res != "" {
			res += "."
		}
		//fmt.Println("read domain name length", length)
		for i := 0; i < int(length); i++ {
			c, err := buffer.ReadU8()
			if err != nil {
				return "", err
			}
			res += string(c)
		}
	}
	if posRestored != 0 {
		buffer.SetPosition(posRestored)
	}
	return res, nil

}
