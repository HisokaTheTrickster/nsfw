package main

import (
	"bytes"
)

type DNS struct {
	Header  DNSHeader
	Query DNSQuery
	Answer  DNSRecord
}

func (d *DNS) ToBytes(recordStat RecordStatus) []byte {

	// Check if the relevent flags are enabled an only convert those struct to bytes
	// For now converting Header and query
	// If domain exists but no record, then ignore the record

	rawData := bytes.Buffer{}

	headerEncoded := d.Header.ToBytesBuffer()
	headerEncoded.WriteTo(&rawData)

	queryEncoded := d.Query.ToBytesBuffer()
	queryEncoded.WriteTo(&rawData)

	if recordStat == RECORD_FOUND_IN_LOCAL {
		recordEncoded := d.Answer.ToBytesBuffer()
		recordEncoded.WriteTo(&rawData)
	}

	return rawData.Bytes()

}

type DNSHeader struct {
	ID                      uint16
	Flags                   uint16
	QuestionCount           uint16
	AnswerCount             uint16
	AuthorityCount          uint16
	AdditionalResourceCount uint16
}

func (d *DNSHeader) ToBytesBuffer() bytes.Buffer {
	encodedMessage := bytes.Buffer{}
	packetBinaryWrite(&encodedMessage, d.ID, d.Flags, d.QuestionCount, d.AnswerCount, d.AuthorityCount, d.AdditionalResourceCount)
	return encodedMessage
}

type DNSQuery struct {
	QueryLabel []string
	QType      uint16
	QClass     uint16
	QueryPtr   uint16
}

func (d *DNSQuery) ToBytesBuffer() bytes.Buffer {

	encodedMessage := bytes.Buffer{}
	for _, label := range d.QueryLabel {
		// first appened the length of the label (1 byte) and then appened the label
		packetBinaryWrite(&encodedMessage, uint8(len(label)), []byte(label))
	}
	// Indicate end of labels
	packetBinaryWrite(&encodedMessage, []byte{0x00})

	packetBinaryWrite(&encodedMessage, d.QType, d.QClass)
	return encodedMessage
}

type DNSRecord struct {
	NamePtr    uint16
	RecordType uint16
	Class      uint16
	TTL        uint32
	RDlength   uint16
	RData      []uint8
}

func (d *DNSRecord) ToBytesBuffer() bytes.Buffer {

	encodedMessage := bytes.Buffer{}
	packetBinaryWrite(&encodedMessage, d.NamePtr, d.RecordType, d.Class, d.TTL, d.RDlength)
	for _, addressOctate := range d.RData {
		packetBinaryWrite(&encodedMessage, addressOctate)
	}

	return encodedMessage

}

type DNSLocalCache struct {
	Rtype uint16 `json:"rtype"`
	Ttl   uint32 `json:"ttl"`
	Value string `json:"value"`
}
