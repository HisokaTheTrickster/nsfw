package utils

import (
	"bytes"
	"encoding/binary"
)

const DEF_DNS_FLAG uint16 = 0

// Main function to handle DNS request

func DNSRequestHandler(buff *bytes.Buffer) (DNS, error) {

	//data := buff.Bytes()
	//fmt.Printf("Bytes Recieved, % x\n", data)

	// Read the flags in the header and identify different fields
	// Exactact the data from different fields

	// If the DNS request is a query, read the different labelds and
	// identify the IP address

	// If its not a DNS request, discard the packet, or just ignore.
	// Create a response packet with the proper headers enabled

	// Send response

	dnsRequest := DNS{}

	// Extract Header : first 12 bytes will be stored in struct DNSHeader / DNS.Header
	headerBytes := buff.Next(12)
	err := binary.Read(bytes.NewBuffer(headerBytes), binary.BigEndian, &dnsRequest.Header)
	if err != nil {
		return dnsRequest, err
	}

	//extractQueries()

	//getResponses()

	//responseToBytes()

	return dnsRequest, err

}
