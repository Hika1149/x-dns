package dns

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
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

	for _, s := range parsedData {

		bs := []byte(s)

		if err := buffer.WriteU8(bs[0]); err != nil {
			return err
		}

	}

	return nil

}
