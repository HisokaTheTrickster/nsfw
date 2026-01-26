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

func FetchRecordFromLocal(localDnsCache map[string][]DNSLocalCache, dnsRequest *DNS ) (RecordStatus, DNSLocalCache) {

	var requestUrl string
	var requestRecordType uint16 = dnsRequest.Query.QType

	// Construct the Request URL
	for _, label := range dnsRequest.Query.QueryLabel {
		requestUrl += label
		requestUrl += "."
	}
	requestUrl = requestUrl[:len(requestUrl)-1]

	// Check if domain exists locally.
	lookupRecords, ifLocalDomainExist := localDnsCache[requestUrl]

	if !ifLocalDomainExist {
		return  DOMAIN_NOT_FOUND_IN_LOCAL, DNSLocalCache{}
	}

	// Check if record exists locally
	for _, rec := range lookupRecords {
		if requestRecordType == rec.Rtype {
			return RECORD_FOUND_IN_LOCAL, rec
		} 
	}
	
	return DOMAIN_EXISTS_NO_RECORD, DNSLocalCache{}

}

func CraftResponseRecord(dnsRespone *DNS, recordStat RecordStatus ,recordFromLocalCache DNSLocalCache) {

	// Since record is not present in local cache, skip creafting Record response
	if recordStat == DOMAIN_EXISTS_NO_RECORD {
		return
	}

	// In the answer header, the name field points to the pointer where the query is requested.
	// This is called query compressions. Insted of the copying the complete query, the answer field just point to it
	queryPointer    := uint16(0xc0) << 8 // This is an indicator. If its 11, then 
	offSet          := uint16(12)
	
	dnsRespone.Answer.NamePtr = queryPointer + offSet
	

	// NEEDS CHANGE: chagne the record type from string to INT for quicker lookup
	switch recordFromLocalCache.Rtype {

	case TypeA:
		ipAddress := []byte(net.ParseIP(recordFromLocalCache.Value))
		dnsRespone.Answer.RecordType = TypeA
		dnsRespone.Answer.RDlength = 4
		dnsRespone.Answer.RData = ipAddress[12:16]

	case TypeAAAA:
		ipAddress := []byte(net.ParseIP(recordFromLocalCache.Value))
		dnsRespone.Answer.RecordType = TypeAAAA
		dnsRespone.Answer.RDlength = 16
		dnsRespone.Answer.RData = ipAddress
	}

	dnsRespone.Answer.Class = 1
	dnsRespone.Answer.TTL = recordFromLocalCache.Ttl

	return
}


func SetHeadersAndFields(dnsRequest, dnsRespone *DNS, recordStat RecordStatus) {

	// Set the headers
	dnsRespone.Header.ID = dnsRequest.Header.ID
	dnsRespone.Header.Flags = dnsRequest.Header.Flags

	dnsRespone.Header.Flags |= 1 << 15            // Response packet
	dnsRespone.Header.Flags |= 1 << 10            // This is the Authority for the domain
	dnsRespone.Header.Flags |= 1 << 7          // Recursion avilable

	dnsRespone.Header.QuestionCount = dnsRequest.Header.QuestionCount
	
	dnsRespone.Header.AuthorityCount = 0
	dnsRespone.Header.AdditionalResourceCount = 0

	if recordStat == RECORD_FOUND_IN_LOCAL {
		dnsRespone.Header.AnswerCount = uint16(1)
	}

	if recordStat == DOMAIN_EXISTS_NO_RECORD {
		dnsRespone.Header.AnswerCount = uint16(0)
		// set the code to nxdomain
		dnsRespone.Header.Flags |= 3
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
