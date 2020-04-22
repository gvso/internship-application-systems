package icmp

import (
	"encoding/binary"
	"errors"
)

// A Message represents an ICMP message.
type Message struct {
	Type     int
	Code     int
	Checksum int
	Body     Body
}

// Body represents the body of a ICMP request.
type Body struct {
	ID   int
	Seq  int
	Data []byte
}

var errTooShort = errors.New("Message is too short")

func checksum(b []byte) uint16 {
	sum := uint16(0)

	// See https://stackoverflow.com/questions/20247551/icmp-echo-checksum.
	for i := 0; i < len(b); i += 2 {
		sum += uint16(b[i+1])<<8 | uint16(b[i])
	}

	return ^sum
}

// marshal encodes the ICMP body into array of bytes
func (eb *Body) marshal() []byte {
	id := uint16(eb.ID)
	seq := uint16(eb.Seq)

	b := []byte{byte(id >> 8), byte(id), byte(seq >> 8), byte(seq)}

	return append(b, eb.Data...)
}

// echoUnmarshal decodes message from array of bytes
func parseBody(b []byte) Body {

	body := Body{
		ID:   int(binary.BigEndian.Uint16(b[0:2])),
		Seq:  int(binary.BigEndian.Uint16(b[2:4])),
		Data: b[4:],
	}

	return body
}

// Marshal encodes the message into array of bytes.
func (m *Message) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}

	b = append(b, m.Body.marshal()...)

	s := checksum(b)

	b[2] = byte(s)
	b[3] = byte(s >> 8)

	return b, nil
}

// ParseMessage decodes message from array of bytes.
func ParseMessage(proto int, b []byte) (*Message, error) {
	// ICMP packets should have at least 4 bytes.
	if len(b) < 4 {
		return nil, errTooShort
	}

	msg := &Message{
		Type:     int(b[0]),
		Code:     int(b[1]),
		Checksum: int(binary.BigEndian.Uint16(b[2:4])),
		Body:     parseBody(b[4:]),
	}

	return msg, nil
}
