package icmp

// NewEchoRequest returns a new message instance for an echo request.
func NewEchoRequest(proto string, b Body) Message {
	return Message{
		Type: 8,
		Code: 0,
		Body: b,
	}
}
