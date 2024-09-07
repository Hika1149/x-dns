package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
)

type DNSPacket struct {
	Header    DNSHeader
	Questions []*DNSQuestion
}

func NewDNSPacket() *DNSPacket {
	return &DNSPacket{
		Header:    *NewDNSHeader(),
		Questions: make([]*DNSQuestion, 0),
	}
}

// FromBuffer reading DNS Info from the buffer  */
func (p *DNSPacket) FromBuffer(buffer buffer.BytePacketBuffer) error {

	err := p.Header.Read(&buffer)
	if err != nil {
		return err
	}

	fmt.Println("Question Count", p.Header.QuestionCount)

	p.Questions = make([]*DNSQuestion, p.Header.QuestionCount)

	for i := 0; i < int(p.Header.QuestionCount); i++ {
		p.Questions[i] = &DNSQuestion{}
		err = p.Questions[i].Read(&buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DNSPacket) Write(buffer *buffer.BytePacketBuffer) error {

	err := p.Header.Write(buffer)
	if err != nil {
		return err
	}

	for _, q := range p.Questions {
		if err := q.Write(buffer); err != nil {
			return err
		}
	}

	return nil

}
