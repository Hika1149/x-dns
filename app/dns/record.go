package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
	"strconv"
	"strings"
)

type Record struct {
	// The domain name to which this resource record pertains
	Name string

	// 1 for A records 5 for CNAME records
	Type uint16

	Class uint16

	TTL uint32

	Length uint16

	Data string
}

func (r *Record) Read(buffer buffer.BufferReader) error {
	var err error
	r.Name, err = DecodeDomainName(buffer)
	if err != nil {
		return err
	}
	r.Type, _ = buffer.ReadU16()
	r.Class, _ = buffer.ReadU16()
	r.TTL, _ = buffer.ReadU32()
	r.Length, _ = buffer.ReadU16()

	// only support A records
	ip := ""
	for i := 0; i < int(r.Length); i++ {
		b, _ := buffer.ReadU8()
		ip += fmt.Sprintf("%d", b)
		if i < int(r.Length)-1 {
			ip += "."

		}

	}
	r.Data = ip
	return nil
}

func (r *Record) Write(buffer buffer.BufferWriter) error {

	domainBytes := encodeDomainName(r.Name)
	for _, b := range domainBytes {
		if err := buffer.WriteU8(b); err != nil {
			return err
		}
	}

	if err := buffer.WriteU16(r.Type); err != nil {
		return err
	}

	if err := buffer.WriteU16(r.Class); err != nil {
		return err
	}

	if err := buffer.WriteU32(r.TTL); err != nil {
		return err
	}
	if err := buffer.WriteU16(r.Length); err != nil {
		return err
	}

	parsedData := strings.Split(r.Data, ".")
	// only support A records
	for _, s := range parsedData {
		intVal, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if err := buffer.WriteU8(byte(intVal)); err != nil {
			return err
		}

	}
	return nil
}

func (r *Record) String() string {
	return fmt.Sprintf("Record{Name: %v, Type: %v, Class: %v, TTL: %v, Length: %v, Data: %v}", r.Name, r.Type, r.Class, r.TTL, r.Length, r.Data)
}
