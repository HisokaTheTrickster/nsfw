package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const DEF_DNS_FLAG uint16 = 0

// Main function to handle DNS request

func DNSRequestHandler(buff *bytes.Buffer) (DNS, error) {

	data := buff.Bytes()
	fmt.Printf("Bytes Recieved, % x\n", data)

	// Read the flags in the header and identify different fields
	// Exactact the data from different fields

	// If the DNS request is a query, read the different labelds and
	// identify the IP address

	// If its not a DNS request, discard the packet, or just ignore.
	// Create a response packet with the proper headers enabled

	// Send response

	dnsRequest := DNS{}

	// Extract Header : first 12 bytes will be stored in struct DNSHeader / DNS.Header
	errHeader := extractHeader(&dnsRequest, buff)
	if errHeader != nil {
		return dnsRequest, errHeader
	}

	// Discard if its a response packet
	if dnsRequest.Header.Flags&0x8000 != 0 {
		return dnsRequest, errors.New("its not a query packet")
	}

	errQuery := extractQueries(&dnsRequest, buff)
	if errQuery != nil {
		return dnsRequest, errQuery
	}

	//getResponses()

	//responseToBytes()

	return dnsRequest, nil

}

func extractHeader(dnsRequest *DNS, buff *bytes.Buffer) error {
	err := binary.Read(bytes.NewBuffer(buff.Next(12)), binary.BigEndian, &dnsRequest.Header)
	return err
}

func extractQueries(dnsRequest *DNS, buff *bytes.Buffer) error {

	noOfQueries := dnsRequest.Header.QuestionCount

	for range noOfQueries {
		dnsQuery := DNSQuery{}
		readLen, err := buff.ReadByte()
		if err != nil {
			return errors.New("unable to read bytes")
		}
		lenOfLabel := int(readLen)

		// extract each labels
		for lenOfLabel != 0 {
			dnsQuery.QueryLabel = append(dnsQuery.QueryLabel, string(buff.Next(lenOfLabel)))
			readLen, err = buff.ReadByte()
			if err != nil {
				return errors.New("unable to read bytes")
			}
			lenOfLabel = int(readLen)
		}
		dnsQuery.QType = binary.BigEndian.Uint16(buff.Next(2))
		dnsQuery.QClass = binary.BigEndian.Uint16(buff.Next(2))

		dnsRequest.Queries = append(dnsRequest.Queries, dnsQuery)
	}

	return nil
}
