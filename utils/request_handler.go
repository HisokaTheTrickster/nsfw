package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

// Main function to handle DNS request

func DNSRequestHandler(buff *bytes.Buffer) error {

	data := buff.Bytes()
	fmt.Printf("Bytes Recieved, % x\n", data)

	// Read the flags in the header and identify different fields
	// Exactact the data from different fields

	// If the DNS request is a query, read the different labelds and
	// identify the IP address

	// If its not a DNS request, discard the packet, or just ignore.
	// Create a response packet with the proper headers enabled

	// Send response

	//extractHeader()

	//extractQueries()

	//getResponses()

	//responseToBytes()

	return nil

}
