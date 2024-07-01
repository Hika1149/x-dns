package dns

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
