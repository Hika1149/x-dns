package dns

import (
	"strings"
)

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func encodeDomainName(domain string) []byte {
	// domain name encoded as a sequence of labels
	// where each label consists of a length octet followed by that
	// many octets
	// the domain name is terminated by a zero length octet
	// the domain name is case-insensitive

	res := make([]byte, 0)

	splitS := strings.Split(domain, ".")

	for _, s := range splitS {

		l := len(s)

		b := []byte(s)

		//fmt.Println("label length", byte(l), "label", b)
		res = append(res, byte(l))
		res = append(res, b...)
	}
	// terminating zero length octet
	res = append(res, 0x00)
	return res

}
