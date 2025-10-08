package utils

// for checking starts of the record. Local/Google etc
type RecordStatus int

const (
	NO_ISSUE                = 0
	RECORD_FOUND_LOCALLY    = 1
	DOMAIN_EXISTS_NO_RECORD = 2
	RECORD_FOUND_REMOTE     = 3
	ERR_REMOTE_DNS_TIMEOUT  = 4
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
