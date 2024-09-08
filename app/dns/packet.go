package dns

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/buffer"
)

type DNSPacket struct {
	Header    DNSHeader
	Questions []*DNSQuestion
	Answers   []*Record
}

func NewDNSPacket() *DNSPacket {
	return &DNSPacket{
		Header:    *NewDNSHeader(),
		Questions: make([]*DNSQuestion, 0),
		Answers:   make([]*Record, 0),
	}
}

// FromBuffer reading DNS Info from the buffer  */
func (p *DNSPacket) FromBuffer(buffer buffer.BytePacketBuffer) error {

	err := p.Header.Read(&buffer)
	if err != nil {
		return err
	}

	p.Questions = make([]*DNSQuestion, p.Header.QuestionCount)

	for i := 0; i < int(p.Header.QuestionCount); i++ {
		p.Questions[i] = &DNSQuestion{}
		err = p.Questions[i].Read(&buffer)
		if err != nil {
			fmt.Printf("read question #%v failed: %v\n", i, err)
			return err
		}
	}

	p.Answers = make([]*Record, p.Header.AnswerCount)
	for i := 0; i < int(p.Header.AnswerCount); i++ {
		p.Answers[i] = &Record{}
		err = p.Answers[i].Read(&buffer)
		if err != nil {
			fmt.Printf("read answer #%v failed: %v\n", i, err)
			return err

		}

	}

	return nil
}

func (p *DNSPacket) AddAnswer(answer *Record) {
	p.Answers = append(p.Answers, answer)
	p.Header.AnswerCount = uint16(len(p.Answers))
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
	for _, answer := range p.Answers {
		if err := answer.Write(buffer); err != nil {
			return err
		}
	}

	return nil
}

func (p *DNSPacket) ToByte() ([]byte, error) {
	buf := buffer.NewBytePacketBuffer()
	err := p.Write(buf)
	if err != nil {
		return nil, err
	}
	return buf.ToByte(), nil

}

func (p *DNSPacket) FromByte(data []byte) *DNSPacket {
	packetBuffer := buffer.BytePacketBuffer{
		Buffer: data,
		Pos:    0,
	}
	err := p.FromBuffer(packetBuffer)
	if err != nil {
		fmt.Println("Failed to parse request:", err)
		return p
	}
	return p

}
