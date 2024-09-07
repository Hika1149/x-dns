package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
)

type DNSHeader struct {
	ID uint16

	// 1 for responses, 0 for queries
	Response bool

	// Operation code
	OpCode uint8

	// Authoritative Answer set to 1
	//if the responding server is an authority for the domain name in question
	AuthoritativeAnswer bool

	TruncatedMsg       bool
	RecursionDesired   bool
	RecursionAvailable bool
	Z                  uint8

	ResCode uint8

	QuestionCount   uint16
	AnswerCount     uint16
	AuthorityCount  uint16
	AdditionalCount uint16
}

func NewDNSHeader() *DNSHeader {
	return &DNSHeader{
		ID:                  0,
		Response:            false,
		OpCode:              0,
		AuthoritativeAnswer: false,
		TruncatedMsg:        false,
		RecursionDesired:    false,
		RecursionAvailable:  false,
		Z:                   0,
		ResCode:             0,
		QuestionCount:       0,
		AnswerCount:         0,
		AuthorityCount:      0,
		AdditionalCount:     0,
	}
}

func (h *DNSHeader) Read(buffer buffer.BufferReader) error {

	var err error

	h.ID, err = buffer.ReadU16()
	if err != nil {
		fmt.Println("error getting packet ID")
		return err
	}

	flags, err := buffer.ReadU16()
	if err != nil {
		return err
	}

	h.Response = (flags & 0x8000) != 0

	h.OpCode = uint8((flags >> 11) & 0x0F)

	// 0000 0100 0000 0000
	h.AuthoritativeAnswer = (flags & 0x0400) != 0

	// 0000 0010 0000 0000
	h.TruncatedMsg = (flags & 0x0200) != 0
	h.RecursionDesired = (flags & 0x0100) != 0
	// 0000 0000 1000 0000
	h.RecursionAvailable = (flags & 0x0080) != 0
	// 0000 0000 0111 0000
	h.Z = uint8(flags & 0x0070)
	//h.Z = uint8((flags >> 4) & 0x07)

	h.ResCode = uint8(flags & 0x000F)

	h.QuestionCount, err = buffer.ReadU16()
	if err != nil {
		return err
	}
	h.AnswerCount, err = buffer.ReadU16()
	if err != nil {
		return err
	}
	h.AuthorityCount, err = buffer.ReadU16()
	if err != nil {
		return err
	}
	h.AdditionalCount, err = buffer.ReadU16()
	if err != nil {
		return err
	}

	return nil

}

func (h *DNSHeader) Write(buffer *buffer.BytePacketBuffer) error {

	if err := buffer.WriteU16(h.ID); err != nil {
		return err
	}

	flagByte := uint16(0x0000)

	// 1 bit
	if h.Response {
		flagByte |= 0x8000
	}
	//4 bit (15-12)
	flagByte |= uint16(h.OpCode) << 11

	//1bit
	if h.AuthoritativeAnswer {
		flagByte |= 0x0400
	}

	//1bit
	if h.TruncatedMsg {
		flagByte |= 0x0200
	}

	//1 bit
	if h.RecursionDesired {
		flagByte |= 0x0100
	}

	// 1bit
	if h.RecursionAvailable {
		flagByte |= 0x0080
	}

	// 3bits (7-5)
	flagByte |= uint16(h.Z) << 4

	// 4bits
	flagByte |= uint16(h.ResCode)

	if err := buffer.WriteU16(flagByte); err != nil {
		return err
	}

	if err := buffer.WriteU16(h.QuestionCount); err != nil {
		return err
	}
	if err := buffer.WriteU16(h.AnswerCount); err != nil {
		return err
	}
	if err := buffer.WriteU16(h.AuthorityCount); err != nil {
		return err
	}
	if err := buffer.WriteU16(h.AdditionalCount); err != nil {
		return err
	}

	return nil
}
