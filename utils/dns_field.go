package utils

import (
	"bytes"
)

type DNS struct {
	Header  DNSHeader
	Queries []DNSQuery
	Answer  []DNSRecords
}

func (d *DNS) ToBytes() []byte {

	// Check if the relevent flags are enabled an only convert those struct to bytes
	// For now converting Header and query

	headerEncoded := d.Header.ToBytes()

	queryEncoded := bytes.Buffer{}
	for i := 0; i < len(d.Queries); i++ {
		ithQuery := d.Queries[i].ToBytes()
		queryEncoded.Write(ithQuery.Bytes())
	}

	recordEncoded := bytes.Buffer{}
	for i := 0; i < len(d.Answer); i++ {
		ithRecord := d.Answer[i].ToBytes()
		recordEncoded.Write(ithRecord.Bytes())
	}

	rawData := bytes.Buffer{}
	rawData.Write(headerEncoded.Bytes())
	rawData.Write(queryEncoded.Bytes())
	rawData.Write(recordEncoded.Bytes())

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

func (d *DNSHeader) ToBytes() bytes.Buffer {
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

func (d *DNSQuery) ToBytes() bytes.Buffer {

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

type DNSRecords struct {
	NamePtr    uint16
	RecordType uint16
	Class      uint16
	TTL        uint32
	RDlength   uint16
	RData      []uint8
}

func (d *DNSRecords) ToBytes() bytes.Buffer {

	encodedMessage := bytes.Buffer{}
	packetBinaryWrite(&encodedMessage, d.NamePtr, d.RecordType, d.Class, d.TTL, d.RDlength)
	for _, addressOctate := range d.RData {
		packetBinaryWrite(&encodedMessage, addressOctate)
	}

	return encodedMessage

}

type DNSdatabase struct {
	Rtype uint16 `json:"rtype"`
	Ttl   uint32 `json:"ttl"`
	Value string `json:"value"`
}
