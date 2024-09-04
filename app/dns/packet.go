package dns

import "github.com/codecrafters-io/dns-server-starter-go/app/buffer"

type DNSPacket struct {
	Header DNSHeader
}

func NewDNSPacket() *DNSPacket {
	return &DNSPacket{
		Header: *NewDNSHeader(),
	}
}

// FromBuffer reading DNS Info from the buffer  */
func (p *DNSPacket) FromBuffer(buffer buffer.BytePacketBuffer) error {

	err := p.Header.Read(&buffer)
	if err != nil {
		return err
	}

	return nil
}

func (p *DNSPacket) Write(buffer *buffer.BytePacketBuffer) error {

	err := p.Header.Write(buffer)
	if err != nil {
		return err
	}

	return nil

}
