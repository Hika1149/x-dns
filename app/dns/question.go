package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
	"runtime/debug"
)

type DNSQuestion struct {
	Name string
	// 2-byte int, the type of record (A, MX, etc)
	QType uint16
	// 2-byte int, the class of record (IN, CS, etc)
	QClass uint16
}

func (q *DNSQuestion) Read(buffer buffer.BufferReader) error {

	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			debug.PrintStack()
		}
	}()

	q.Name, err = q.ReadeDomainName(buffer)
	if err != nil {
		fmt.Println("read domain name failed", err)
		return err
	}
	q.QType, err = buffer.ReadU16()
	if err != nil {
		fmt.Println("read QType failed", err)
		return err
	}

	q.QClass, err = buffer.ReadU16()
	if err != nil {
		fmt.Println("read QClass failed", err)
		return err
	}

	return nil

}

func (q *DNSQuestion) Write(buffer buffer.BufferWriter) error {

	domainBytes := encodeDomainName(q.Name)

	for _, b := range domainBytes {
		if err := buffer.WriteU8(b); err != nil {
			return err
		}
	}

	if err := buffer.WriteU16(q.QType); err != nil {
		return err
	}

	if err := buffer.WriteU16(q.QClass); err != nil {
		return err
	}
	return nil

}

func (q *DNSQuestion) ReadeDomainName(buffer buffer.BufferReader) (string, error) {

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
