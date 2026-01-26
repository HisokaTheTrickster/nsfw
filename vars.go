package main

import "errors"

// for checking starts of the record. Local/Google etc
type RecordStatus int

const (
	NO_ISSUE                	= RecordStatus(0)
	DOMAIN_NOT_FOUND_IN_LOCAL 	= RecordStatus(1)
	RECORD_FOUND_IN_LOCAL 	 	= RecordStatus(2)
	DOMAIN_EXISTS_NO_RECORD 	= RecordStatus(3)
 	ERR_REMOTE_DNS_TIMEOUT  	= RecordStatus(4)
)

const (
	DNS_DB_PATH      = "records.json"
	DNS_ADDRESS_PORT = ":53"
)

const (
	TypeA     uint16 = 1
	TypeAAAA  uint16 = 28
	TypeMX    uint16 = 15
	TypeCNAME uint16 = 5
)

var (
	errExtractDNSHeader = errors.New("Unable to extract the DNS request Header")
	errExtractDNSQuery = errors.New("Unable to extract the DNS request Query")
)
