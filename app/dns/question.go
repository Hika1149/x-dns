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

	q.Name, err = DecodeDomainName(buffer)
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
