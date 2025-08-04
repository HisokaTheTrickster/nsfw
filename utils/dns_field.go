package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DNS struct {
	// should query be made as a slice????
	// Query []DNSQueries
	Header DNSHeader
	Query  DNSQueries
}

func (d *DNS) ToBytes() bytes.Buffer {

	// Check if the relevent flags are enabled an only convert those struct to bytes
	// For now converting Header and query
	headerEncoded := d.Header.ToBytes()
	queryEncoded := d.Query.ToBytes()

	rawData := bytes.Buffer{}

	rawData.Write(headerEncoded.Bytes())
	rawData.Write(queryEncoded.Bytes())

	return rawData

}

type DNSHeader struct {
	ID                      uint16
	Flags                   uint16
	QuestionCount           uint16
	AnswerCount             uint16
	AuthorityCount          uint16
	AdditionalResourceCount uint16
}

func (d *DNSHeader) ToBytes() bytes.Buffer {
	encodedMessage := bytes.Buffer{}
	packetBinaryWrite(&encodedMessage, d.ID, d.Flags, d.QuestionCount, d.AnswerCount, d.AuthorityCount, d.AdditionalResourceCount)
	return encodedMessage
}

type DNSQueries struct {
	QueryLabel []string
	QType      uint16
	QClass     uint16
}

func (d *DNSQueries) ToBytes() bytes.Buffer {

	encodedMessage := bytes.Buffer{}
	for _, label := range d.QueryLabel {
		packetBinaryWrite(&encodedMessage, len(label), label)
	}
	packetBinaryWrite(&encodedMessage, d.QType, d.QClass)
	return encodedMessage
}

func packetBinaryWrite(buff io.Writer, data ...any) {
	for _, iData := range data {
		binary.Write(buff, binary.BigEndian, iData)
	}
}
