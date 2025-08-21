package utils

import (
	"encoding/json"
	"os"
)

const (
	DNS_DB_PATH = "records.json"
)

func LoadAllDNSCache() (map[string]DNSdatabase, error) {

	// dnsCache := Loading Cache
	// DNS cache is a map of URL and Records

	db := []DNSdatabase{}
	data, err := os.ReadFile(DNS_DB_PATH)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	dnsCache := make(map[string]DNSdatabase)
	for _, record := range db {
		dnsCache[record.Name] = record
	}

	return dnsCache, nil

}
