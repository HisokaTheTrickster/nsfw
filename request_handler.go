package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const DEF_DNS_FLAG uint16 = 0

// Main function to handle DNS request

func ExtractRequest(buff *bytes.Buffer) (DNS, error) {

	dnsRequest := DNS{}
	var err error

	err = extractHeader(&dnsRequest, buff)
	if err != nil {
		return dnsRequest, err
	}

	err = extractQueries(&dnsRequest, buff)
	if err != nil {
		return dnsRequest, err
	}

	return dnsRequest, nil
}

func extractHeader(dnsRequest *DNS, buff *bytes.Buffer) error {
	// Extract Header : first 12 bytes will be stored in struct DNSHeader / DNS.Header
	err := binary.Read(bytes.NewBuffer(buff.Next(12)), binary.BigEndian, &dnsRequest.Header)

	if err != nil {
		return errExtractDNSHeader
	}

	return nil

}

func extractQueries(dnsRequest *DNS, buff *bytes.Buffer) error {

	readLen, err := buff.ReadByte()

	if err != nil {
		return errExtractDNSQuery
	}
	lenOfLabel := int(readLen)

	// On the wire, first comes the length of label, followed by the label itself.
	// When length of label becomes, 0 we know that there are no more labels
	for lenOfLabel != 0 {
		dnsRequest.Query.QueryLabel = append(dnsRequest.Query.QueryLabel, string(buff.Next(lenOfLabel)))
		readLen, err = buff.ReadByte()
		if err != nil {
			return errExtractDNSQuery
		}
		lenOfLabel = int(readLen)
	}

	dnsRequest.Query.QType = binary.BigEndian.Uint16(buff.Next(2))
	dnsRequest.Query.QClass = binary.BigEndian.Uint16(buff.Next(2))

	return nil
}

func CraftResponseRecord(dnsRespone *DNS, recordStat RecordStatus, recordFromLocalCache DNSLocalCache) {

	// Since record is not present in local cache, skip creafting Record response
	if recordStat == DOMAIN_EXISTS_NO_RECORD || recordStat == DOMAIN_BLOCKED {
		return
	}

	// In the answer header, the name field points to the pointer where the query is requested.
	// This is called query compressions. Insted of the copying the complete query, the answer field just point to it
	queryPointer := uint16(0xc0) << 8 // This is an indicator. If its 11, then
	offSet := uint16(12)

	dnsRespone.Answer.NamePtr = queryPointer + offSet

	// NEEDS CHANGE: chagne the record type from string to INT for quicker lookup
	switch recordFromLocalCache.Rtype {

	case TypeA:
		ipAddress := net.ParseIP(recordFromLocalCache.Value)
		if ipAddress == nil {
			log.Printf("Invalid IPv4 address in local cache: %s", recordFromLocalCache.Value)
			return
		}
		dnsRespone.Answer.RecordType = TypeA
		dnsRespone.Answer.RDlength = 4
		dnsRespone.Answer.RData = []byte(ipAddress.To4())

	case TypeAAAA:
		ipAddress := net.ParseIP(recordFromLocalCache.Value)
		if ipAddress == nil {
			log.Printf("Invalid IPv6 address in local cache: %s", recordFromLocalCache.Value)
			return
		}
		dnsRespone.Answer.RecordType = TypeAAAA
		dnsRespone.Answer.RDlength = 16
		dnsRespone.Answer.RData = []byte(ipAddress.To16())
	}

	dnsRespone.Answer.Class = 1
	dnsRespone.Answer.TTL = recordFromLocalCache.Ttl
}

func SetHeadersAndFields(dnsRequest, dnsRespone *DNS, recordStat RecordStatus) {

	// Set the headers
	dnsRespone.Header.ID = dnsRequest.Header.ID
	dnsRespone.Header.Flags = dnsRequest.Header.Flags

	dnsRespone.Header.Flags |= 1 << 15 // Response packet
	dnsRespone.Header.Flags |= 1 << 10 // This is the Authority for the domain
	dnsRespone.Header.Flags |= 1 << 7  // Recursion avilable

	dnsRespone.Header.QuestionCount = dnsRequest.Header.QuestionCount

	dnsRespone.Header.AuthorityCount = 0
	dnsRespone.Header.AdditionalResourceCount = 0

	switch recordStat {
	case RECORD_FOUND_IN_LOCAL:
		dnsRespone.Header.AnswerCount = 1
		// RCODE = NOERROR (0) → default

	case DOMAIN_EXISTS_NO_RECORD:
		dnsRespone.Header.AnswerCount = 0
		// RCODE = NOERROR (0) → domain exists, no record of this type

	case DOMAIN_BLOCKED:
		dnsRespone.Header.AnswerCount = 0
		// Set RCODE = NXDOMAIN (3)
		dnsRespone.Header.Flags &^= 0xF // clear last 4 bits (existing RCODE)
		dnsRespone.Header.Flags |= 0x3  // set NXDOMAIN
	}

	dnsRespone.Query = dnsRequest.Query

}

func DiscardRequest(dnsRequest *DNS) bool {

	// check if the DNS flag for request is set to true
	if dnsRequest.Header.Flags&0x8000 != 0 {
		log.Println("Not a DNS request. Dropping the packet")
		return true
	}

	return false

}
