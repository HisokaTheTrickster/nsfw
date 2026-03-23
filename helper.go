package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func LoadLocalDNSRecords(path string) (map[string][]DNSLocalCache, error) {

	db := make(map[string][]DNSLocalCache)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("no local records found")
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, errors.New("invalid record format, skipping records.json")
	}

	return db, nil
}

// normalizeDomain converts a domain to lowercase and removes trailing dots
func normalizeDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimSpace(domain)
	domain = strings.TrimSuffix(domain, ".")
	return domain
}

// parseLine extracts a domain from a hosts/adblock style line
func parseLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return ""
	}

	// Hosts format: "0.0.0.0 domain.com"
	fields := strings.Fields(line)
	if len(fields) >= 2 && (fields[0] == "0.0.0.0" || fields[0] == "127.0.0.1") {
		return normalizeDomain(fields[1])
	}

	// Adblock style: "||domain.com^" or "||domain.com^$options"
	if domain, ok := strings.CutPrefix(line, "||"); ok {
		domain, _, _ = strings.Cut(domain, "^")
		return normalizeDomain(domain)
	}

	// Plain domain (if nothing else)
	return normalizeDomain(line)
}

// loadBlocklist reads the file and returns a map of blocked domains
func loadBlocklist(filePath string) (map[string]struct{}, error) {
	blocked := make(map[string]struct{}, 150000) // preallocate for 130k entries

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := parseLine(scanner.Text())
		if domain != "" {
			blocked[domain] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return blocked, nil
}

func packetBinaryWrite(buff io.Writer, data ...any) {
	// writing data to bytes
	for _, iData := range data {
		binary.Write(buff, binary.BigEndian, iData)
	}
}

func FetchFromPublicDNS(inputBuffer []byte, inputBuffSize int) ([]byte, error) {

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

	_, err = conn.Write(inputBuffer[:inputBuffSize])
	if err != nil {
		return dnsServerBuffer, err
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := conn.ReadFromUDP(dnsServerBuffer)
	if err != nil {
		return dnsServerBuffer, err
	}

	return dnsServerBuffer[:n], err

}

func FetchRecordFromLocal(localDnsCache map[string][]DNSLocalCache, localBlockDomain map[string]struct{}, dnsRequest *DNS) (RecordStatus, DNSLocalCache) {

	var requestUrl string
	var requestRecordType uint16 = dnsRequest.Query.QType

	// Guard against malformed requests with no labels
	if len(dnsRequest.Query.QueryLabel) == 0 {
		return DOMAIN_NOT_FOUND_IN_LOCAL, DNSLocalCache{}
	}

	// Construct the Request URL
	for _, label := range dnsRequest.Query.QueryLabel {
		requestUrl += label
		requestUrl += "."
	}
	requestUrl = requestUrl[:len(requestUrl)-1]

	// Local cache takes priority — if a matching record exists, serve it
	lookupRecords, ifLocalDomainExist := localDnsCache[requestUrl]
	if ifLocalDomainExist {
		for _, rec := range lookupRecords {
			if requestRecordType == rec.Rtype {
				return RECORD_FOUND_IN_LOCAL, rec
			}
		}
	}

	// Blocklist applies to everything not directly served from local cache
	if _, blocked := localBlockDomain[requestUrl]; blocked {
		return DOMAIN_BLOCKED, DNSLocalCache{}
	}

	// Domain exists locally but no matching record type
	if ifLocalDomainExist {
		return DOMAIN_EXISTS_NO_RECORD, DNSLocalCache{}
	}

	return DOMAIN_NOT_FOUND_IN_LOCAL, DNSLocalCache{}

}
