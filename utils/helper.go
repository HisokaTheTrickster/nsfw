package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func LoadAllDNSCache() (map[string][]DNSdatabase, error) {

	db := make(map[string][]DNSdatabase)

	data, err := os.ReadFile("records.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	fmt.Println(db)

	return db, err

}

func packetBinaryWrite(buff io.Writer, data ...any) {
	// writing data to bytes
	for _, iData := range data {
		binary.Write(buff, binary.BigEndian, iData)
	}
}

func CraftResponse(lookup DNSdatabase, record *DNSRecords) {

	switch lookup.Rtype {

	// NEEDS CHANGE: chagne the record type from string to INT for quicker lookup

	case TypeA:
		ipAddress := []byte(net.ParseIP(lookup.Value))
		record.RecordType = TypeA
		record.RDlength = 4
		record.RData = ipAddress[12:16]

	case TypeAAAA:
		ipAddress := []byte(net.ParseIP(lookup.Value))
		record.RecordType = TypeAAAA
		record.RDlength = 16
		record.RData = ipAddress
	}

	record.Class = 1
	record.TTL = lookup.Ttl

}

func FetchFromNet(inputBufferPtr *[]byte, inputBuffSize int) ([]byte, error) {

	dnsServerBuffer := make([]byte, 512)

	serverAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		return dnsServerBuffer, err
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return dnsServerBuffer, err
	}
	defer conn.Close()

	_, err = conn.Write((*inputBufferPtr)[:inputBuffSize])
	if err != nil {
		return dnsServerBuffer, err
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := conn.ReadFromUDP(dnsServerBuffer)
	_ = n

	if err != nil {
		fmt.Println("Error or timeout:", err)
		return dnsServerBuffer, err
	}

	return dnsServerBuffer[:n], err

}

/*
func FetchFromNet(requestUrl string, requestType uint16) (DNSdatabase, bool) {

	recordToSave := DNSdatabase{}
	ifExist := false

	switch requestType {

	case TypeA:
		recordToSave.Rtype = TypeA
		ips, err := net.LookupIP(requestUrl)

		if err == nil {
			ifExist = true
			for _, ip := range ips {
				if ip.To4() != nil {
					recordToSave.Value = ip.String()
				}
			}
		}

	case TypeAAAA:
		recordToSave.Rtype = TypeAAAA
		ips, err := net.LookupIP(requestUrl)

		if err == nil {
			ifExist = true

			for _, ip := range ips {
				if ip.To4() == nil {
					recordToSave.Value = ip.String()
				}
			}
		}

	// CNAME record
	case TypeCNAME:
		recordToSave.Rtype = TypeCNAME
		cname, err := net.LookupCNAME(requestUrl)
		if err == nil {
			ifExist = true
		}
		recordToSave.Value = cname

	}

	return recordToSave, ifExist

}
*/
