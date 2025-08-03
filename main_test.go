package main

import (
	"bytes"
	"nsfw/utils"
	"testing"
)

func TestDNSResolver(t *testing.T) {

	// Create a dummy header
	testHeader := utils.DNSHeader{ID: 23, Flags: utils.DEF_DNS_FLAG, QuestionCount: 1, AnswerCount: 0, AuthorityCount: 0, AdditionalResourceCount: 0}
	testQuery := utils.DNSQueries{QueryLabel: []string{"google", "com"}, QType: 1, QClass: 1}

	rawHeader := testHeader.ToBytes()
	rawQuery := testQuery.ToBytes()

	finalRaw := bytes.Buffer{}
	finalRaw.Write(rawHeader.Bytes())
	finalRaw.Write(rawQuery.Bytes())

	err := utils.DNSRequestHandler(&finalRaw)

	if err != nil {
		t.Errorf("Error occured")
	}

}
