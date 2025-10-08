package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
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

func FetchRecord(allDnsCache map[string][]DNSdatabase, dnsRequest *DNS, dnsResponse *DNS, inputBufferPtr *[]byte, inputBuffSize int) (RecordStatus, []byte) {

	log.Println("Checking Database for a response")

	// query pointer to record mapping
	// In the answer header, the name field points to the pointer where the query is requested.
	// This is query compressions. Insted of the complete query, you just point to it

	// extract the IP for each query
	// Right now, we are only resolving a single record.

	var (
		pointer         = uint16(0xc0) << 8
		offSet          = uint16(12)
		recStat         = RecordStatus(NO_ISSUE)
		dnsServerBuffer = make([]byte, 1024)
	)

	var err error

	// There could be several queries in a single DNS packet. We itrate throught each of them
	// For now, we only are only checking the first record.

	for _, query := range dnsRequest.Queries {

		responseRecord := DNSRecords{}
		responseRecord.NamePtr = pointer + offSet

		var requestUrl string
		var requestType uint16 = query.QType

		for index, label := range query.QueryLabel {
			requestUrl += label
			if index != len(query.QueryLabel)-1 {
				requestUrl += "."
			}
			offSet += uint16(1 + len(label))
		}
		// add offset for qtype and qclass
		offSet += uint16(4)

		ifLocalDomainExist, ifLocalRecExist := false, false

		// Check if record exist locally. If not, send it to Google DNS

		// If domain exists, and no record is found, just send an empty response.
		// If domain does not exist, send NXDOMAIN

		localLookup, ifLocalDomainExist := allDnsCache[requestUrl]

		if ifLocalDomainExist {

			for _, rec := range localLookup {
				if requestType == rec.Rtype {
					ifLocalRecExist = true
					recStat = RECORD_FOUND_LOCALLY
					CraftResponse(rec, &responseRecord)
					break
				}
			}

			if !ifLocalRecExist {
				// send no error and no response
				recStat = DOMAIN_EXISTS_NO_RECORD

			}

		} else {
			dnsServerBuffer, err = FetchFromNet(inputBufferPtr, inputBuffSize)
			if err != nil {
				recStat = ERR_REMOTE_DNS_TIMEOUT
			} else {
				recStat = RECORD_FOUND_REMOTE
			}

		}

		if ifLocalRecExist {
			dnsResponse.Answer = append(dnsResponse.Answer, responseRecord)
		}

		// We are only checking the first query. Not all of them
		break

	}

	/*
		if len(dnsResponse.Answer) == 0 {
			return RECORD_NOT_FOUND
		}
	*/

	return recStat, dnsServerBuffer

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

	dnsRespone.Queries = dnsRequest.Queries

}

func DiscardRequest(dnsRequest *DNS) bool {

	if dnsRequest.Header.Flags&0x8000 != 0 {
		log.Println("This is not a DNS request. Dropping the packet")
		return true
	}

	return false

}

func NoRecordFound(dnsRequest, dnsRespone *DNS) {

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

	dnsRespone.Queries = dnsRequest.Queries

}
