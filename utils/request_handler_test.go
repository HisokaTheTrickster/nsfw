package utils

import (
	"reflect"
	"testing"
)

func TestCheckHeader(t *testing.T) {

	// Create a dummy header
	t.Run("Check for errors", func(t *testing.T) {
		testDNS := DNS{
			Header: DNSHeader{ID: 23, Flags: DEF_DNS_FLAG, QuestionCount: 1, AnswerCount: 0, AuthorityCount: 0, AdditionalResourceCount: 0},
			Query:  DNSQueries{QueryLabel: []string{"google", "com"}, QType: 1, QClass: 1},
		}

		rawBytes := testDNS.ToBytes()
		_, err := DNSRequestHandler(&rawBytes)

		if err != nil {
			t.Errorf("Error occured")
		}
	})

	t.Run("Check if header is decoded properly", func(t *testing.T) {
		testDNSHeader := DNSHeader{ID: 23, Flags: DEF_DNS_FLAG, QuestionCount: 1, AnswerCount: 0, AuthorityCount: 0, AdditionalResourceCount: 0}
		rawData := testDNSHeader.ToBytes()

		got, _ := DNSRequestHandler(&rawData)

		if !reflect.DeepEqual(got.Header, testDNSHeader) {
			t.Errorf("%v, %v", got, testDNSHeader)
		}

	})

}
