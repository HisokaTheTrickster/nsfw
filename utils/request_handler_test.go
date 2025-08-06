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
			Queries: []DNSQuery{
				{QueryLabel: []string{"google", "com"}, QType: 1, QClass: 1},
			},
		}

		rawBytes := testDNS.ToBytes()
		resultDNS, err := DNSRequestHandler(&rawBytes)

		if err != nil {
			t.Errorf("Error occured")
		}

		if !reflect.DeepEqual(resultDNS, testDNS) {
			t.Errorf("%v, %v", resultDNS, testDNS)
		}
	})

	t.Run("Check if header is decoded properly", func(t *testing.T) {
		// Create two DNS instance and implement only the Header.ToBytes() method
		// The first instance will have the test data,
		// The test will extract the bytes from first instance which will be used as raw packet input to the second instance

		testDNSWant := DNS{
			Header: DNSHeader{ID: 23, Flags: DEF_DNS_FLAG, QuestionCount: 1, AnswerCount: 0, AuthorityCount: 0, AdditionalResourceCount: 0},
		}
		rawData := testDNSWant.Header.ToBytes()

		testDNSGot := DNS{}
		err := extractHeader(&testDNSGot, &rawData)

		if err != nil {
			t.Errorf("%s", err)
		}

		if !reflect.DeepEqual(testDNSWant.Header, testDNSGot.Header) {
			t.Errorf("%v, %v", testDNSWant, testDNSGot)
		}

	})

}
