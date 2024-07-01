package buffer

import (
	"encoding/binary"
	"errors"
)

type BufferWriter interface {
	WriteU8(val uint8) error
	WriteU16(val uint16) error
	WriteU32(val uint32) error
}
type BufferReader interface {
	ReadU8() (uint8, error)
	ReadU16() (uint16, error)
	ReadU32() (uint32, error)
}

type BytePacketBuffer struct {
	Buffer []byte
	Pos    uint16
}

func NewBytePacketBuffer() *BytePacketBuffer {
	//Conventionally, DNS packets are sent using UDP transport and are limited to 512 bytes
	return &BytePacketBuffer{
		Buffer: make([]byte, 512),
		Pos:    0,
	}
}

func (b *BytePacketBuffer) Position() uint16 {
	return b.Pos
}

func (b *BytePacketBuffer) Get(pos uint16) (uint8, error) {
	if err := b.checkBounds(pos, 0); err != nil {
		return 0, err
	}
	return b.Buffer[pos], nil
}

func (b *BytePacketBuffer) GetRange(start, length uint16) ([]byte, error) {
	if err := b.checkBounds(start, length); err != nil {
		return nil, err
	}
	return b.Buffer[start : start+length], nil
}

func (b *BytePacketBuffer) ReadU8() (uint8, error) {
	if err := b.checkBounds(b.Pos, 1); err != nil {
		return 0, err
	}
	val := b.Buffer[b.Pos]
	b.Pos++
	return val, nil
}

func (b *BytePacketBuffer) ReadU16() (uint16, error) {
	if err := b.checkBounds(b.Pos, 2); err != nil {
		return 0, err
	}
	val := binary.BigEndian.Uint16(b.Buffer[b.Pos:])
	b.Pos += 2
	return val, nil
}
func (b *BytePacketBuffer) ReadU32() (uint32, error) {
	if err := b.checkBounds(b.Pos, 4); err != nil {
		return 0, err
	}
	val := binary.BigEndian.Uint32(b.Buffer[b.Pos:])
	b.Pos += 4
	return val, nil

}

func (b *BytePacketBuffer) WriteU8(val uint8) error {
	if err := b.checkBounds(b.Pos, 1); err != nil {
		return err
	}
	b.Buffer[b.Pos] = val
	b.Pos++
	return nil
}

func (b *BytePacketBuffer) WriteU16(val uint16) error {
	if err := b.checkBounds(b.Pos, 2); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(b.Buffer[b.Pos:], val)
	b.Pos += 2
	return nil

}
func (b *BytePacketBuffer) WriteU32(val uint32) error {
	if err := b.checkBounds(b.Pos, 4); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(b.Buffer[b.Pos:], val)
	b.Pos += 4
	return nil

}

func (b *BytePacketBuffer) checkBounds(pos, length uint16) error {
	bufSize := uint16(len(b.Buffer))
	if pos >= bufSize || pos+length > bufSize {
		return errors.New("out of bounds")
	}
	return nil
}
