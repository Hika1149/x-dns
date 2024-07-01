package main

type DNSMessage struct {
	Header   [12]byte
	Question []byte
	//Answer   []byte
}

//func (msg DNSMessage) ToBytes() []byte {
//	result := make([]byte, 0)
//	result = append(result, msg.Header[:]...)
//	result = append(result, msg.Question...)
//	//result = append(result, msg.Answer...)
//	return result
//}

func (msg DNSMessage) ToBytes() []byte {
	lenHeader := len(msg.Header)
	buf := make([]byte, lenHeader+len(msg.Question))
	copy(buf, msg.Header[:])
	copy(buf[lenHeader:], msg.Question[:])
	return buf
}
