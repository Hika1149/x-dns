package udp

import (
	"bufio"
	"net"
)

func Dial(ip string, contents []byte) ([]byte, error) {

	// use dig to test dns server locally

	p := make([]byte, 2048)

	conn, err := net.Dial("udp", ip)
	if err != nil {
		return p, err
	}
	_, err = conn.Write(contents)
	if err != nil {
		return p, err
	}

	size, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		return p, err
	}

	return p[:size], nil

}
