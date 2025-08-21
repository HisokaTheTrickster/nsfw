package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

const DEF_DNS_FLAG uint16 = 0

// Main function to handle DNS request

func ExtractRequest(buff *bytes.Buffer) (DNS, error) {

	dnsRequest := DNS{}

	errHeader := extractHeader(&dnsRequest, buff)
	if errHeader != nil {
		return dnsRequest, errHeader
	}

	errQuery := extractQueries(&dnsRequest, buff)
	if errQuery != nil {
		return dnsRequest, errQuery
	}

	return dnsRequest, nil
}

func extractHeader(dnsRequest *DNS, buff *bytes.Buffer) error {
	// Extract Header : first 12 bytes will be stored in struct DNSHeader / DNS.Header
	err := binary.Read(bytes.NewBuffer(buff.Next(12)), binary.BigEndian, &dnsRequest.Header)

	if err != nil {
		return errors.New("unable to extract dns header")
	}

	return nil

}

func extractQueries(dnsRequest *DNS, buff *bytes.Buffer) error {
	noOfQueries := dnsRequest.Header.QuestionCount

	for range noOfQueries {
		dnsQuery := DNSQuery{}
		readLen, err := buff.ReadByte()

		if err != nil {
			return errors.New("unable to extract dns query")
		}
		lenOfLabel := int(readLen)

		// extract each labels
		for lenOfLabel != 0 {
			dnsQuery.QueryLabel = append(dnsQuery.QueryLabel, string(buff.Next(lenOfLabel)))
			readLen, err = buff.ReadByte()
			if err != nil {
				return errors.New("unable to extract dns query")
			}
			lenOfLabel = int(readLen)
		}

		dnsQuery.QType = binary.BigEndian.Uint16(buff.Next(2))
		dnsQuery.QClass = binary.BigEndian.Uint16(buff.Next(2))

		dnsRequest.Queries = append(dnsRequest.Queries, dnsQuery)

	}

	return nil
}

func FetchFromLocalRecord(allDnsCache map[string]DNSdatabase, dnsRequest *DNS, dnsResponse *DNS) error {

	log.Println("Checking Database for a response")

	// query pointer to record mapping
	// In the answer header, the name field points to the pointer where the query is requested.
	// This is query compressions. Insted of the complete query, you just point to it

	// extract the IP for each query
	pointer := uint16(0xc0) << 8
	offSet := uint16(12)

	for _, query := range dnsRequest.Queries {

		// skip if the query is of type AAAA
		if query.QType == 28 {
			continue
		}

		responseRecord := DNSRecords{}
		responseRecord.NamePtr = pointer + offSet

		var requestUrl string
		for index, label := range query.QueryLabel {
			requestUrl += label
			if index != len(query.QueryLabel)-1 {
				requestUrl += "."
			}
			offSet += uint16(1 + len(label))
		}

		// add offset for qtype and qclass
		offSet += uint16(4)

		record, ifExist := allDnsCache[requestUrl]
		fmt.Println(allDnsCache)
		fmt.Println(requestUrl)

		// check if record exists
		if !ifExist {
			return ErrNoLocalRecord
		}

		ipAddress := []byte(net.ParseIP(record.Address))

		if record.V6 {
			responseRecord.RecordType = 28
			responseRecord.RDlength = 16
			responseRecord.RData = ipAddress

		} else {
			responseRecord.RecordType = 1
			responseRecord.RDlength = 4
			responseRecord.RData = ipAddress[12:16]
		}

		responseRecord.Class = 1
		responseRecord.TTL = 150

		dnsResponse.Answer = append(dnsResponse.Answer, responseRecord)

	}

	if len(dnsResponse.Answer) == 0 {
		return ErrNoLocalRecord
	}

	return nil

}

func CopyRequiredFields(dnsRequest, dnsRespone *DNS) {

	// Set the headers
	dnsRespone.Header.ID = dnsRequest.Header.ID
	dnsRespone.Header.Flags = dnsRequest.Header.Flags

	dnsRespone.Header.Flags |= 1 << 15            // Response packet
	dnsRespone.Header.Flags |= 1 << 10            // This is the Authority for the domain
	dnsRespone.Header.Flags &= 0b1111111011111111 // The packet is not truncated
	dnsRespone.Header.Flags &= 0b1111111110111111 // Recursion not avilable

	dnsRespone.Header.QuestionCount = dnsRequest.Header.QuestionCount
	dnsRespone.Header.AnswerCount = uint16(len(dnsRespone.Answer))

	dnsRespone.Header.AuthorityCount = 0
	dnsRespone.Header.AdditionalResourceCount = 0

	// Set the query section
	dnsRespone.Queries = dnsRequest.Queries

	// response is aldready set

	// convert the headers to bytes.

}

func DiscardRequest(dnsRequest *DNS) bool {

	if dnsRequest.Header.Flags&0x8000 != 0 {
		log.Println("This is not a DNS request. Dropping the packet")
		return true
	}

	return false

}
