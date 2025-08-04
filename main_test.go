package main

import (
	"nsfw/utils"
	"testing"
)

func TestDNSResolver(t *testing.T) {

	// Create a dummy header
	testDNS := utils.DNS{
		Header: utils.DNSHeader{ID: 23, Flags: utils.DEF_DNS_FLAG, QuestionCount: 1, AnswerCount: 0, AuthorityCount: 0, AdditionalResourceCount: 0},
		Query:  utils.DNSQueries{QueryLabel: []string{"google", "com"}, QType: 1, QClass: 1},
	}

	rawBytes := testDNS.ToBytes()
	err := utils.DNSRequestHandler(&rawBytes)

	if err != nil {
		t.Errorf("Error occured")
	}

}
