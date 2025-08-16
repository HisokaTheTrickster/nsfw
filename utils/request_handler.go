package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const DEF_DNS_FLAG uint16 = 0

// Main function to handle DNS request

func RequestHandler(buff *bytes.Buffer) (DNS, error) {

	//data := buff.Bytes()
	//fmt.Printf("Bytes Recieved, % x\n", data)

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

	// getResponses()
	// Future releases

	// responseToBytes()
	// Future releases

	// Discard if its a response packet
	//if dnsRequest.Header.Flags&0x8000 != 0 {
	//	return dnsRequest, errors.New("its not a query packet")
	//}

}

func extractHeader(dnsRequest *DNS, buff *bytes.Buffer) error {
	// Extract Header : first 12 bytes will be stored in struct DNSHeader / DNS.Header
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

func FetchRecord(allDnsCache map[string]DNSdatabase, dnsPacket *DNS) error {

	// query pointer to record mapping
	// In the answer header, the name field points to the pointer where the query is requested.
	// This is query compressions. Insted of the complete query, you just point to it

	if dnsPacket.Header.Flags&0x8000 != 0 {
		return errors.New("this is not a request packet")
	}

	// extract the IP for each query
	pointer := uint16(0xc0) << 8
	offSet := uint16(12)

	for _, query := range dnsPacket.Queries {

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

		if !ifExist {
			fmt.Printf("Skipping query %s as no records exists locally\n", requestUrl)
			continue
			// need to call Google DNS in case the record is not there locally
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

		dnsPacket.Answer = append(dnsPacket.Answer, responseRecord)

	}

	return nil

}
