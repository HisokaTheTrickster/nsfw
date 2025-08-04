package utils

import (
	"bytes"
	"fmt"
)

const DEF_DNS_FLAG uint16 = 0

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
