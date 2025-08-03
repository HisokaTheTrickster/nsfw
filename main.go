package main

import (
	"bytes"
	"encoding/binary"
)

const DEF_DNS_FLAG uint16 = 0

type DNSHeader struct {
	ID                      uint16
	Flags                   uint16
	QuestionCount           uint16
	AnswerCount             uint16
	AuthorityCount          uint16
	AdditionalResourceCount uint16
}

func (d *DNSHeader) ToBytes() []byte {

	encodedMessage := &bytes.Buffer{}
	binary.Write(encodedMessage, binary.BigEndian, d.ID)
	binary.Write(encodedMessage, binary.BigEndian, d.Flags)
	binary.Write(encodedMessage, binary.BigEndian, d.QuestionCount)
	binary.Write(encodedMessage, binary.BigEndian, d.AnswerCount)
	binary.Write(encodedMessage, binary.BigEndian, d.AuthorityCount)
	binary.Write(encodedMessage, binary.BigEndian, d.AdditionalResourceCount)

	return encodedMessage.Bytes()
}
