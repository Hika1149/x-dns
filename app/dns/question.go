package dns

type DNSQuestion struct {
	Name string
	// 2-byte int, the type of record (A, MX, etc)
	QType uint16
	// 2-byte int, the class of record (IN, CS, etc)
	QClass uint16
}
