package main

import (
	"bytes"
	"log"
	"net"
)

type jobRequest struct {
	addr *net.UDPAddr
	data []byte
	size int
}

func main() {

	log.Printf(BANNER)

	// Listen on port 53
	serverAddr, err := net.ResolveUDPAddr("udp", DNS_ADDRESS_PORT)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	// Load records from records.json
	log.Println("loading local records")
	localDnsCache, err := LoadLocalDNSRecords(DNS_DB_PATH)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("local records loaded")
	}

	// Load Domains that needs to be blocked
	log.Println("loading domains that needs to be blocked")
	localBlockDomain, err := loadBlocklist(DOMAIN_BLOCK_LIST)
	if err != nil {
		log.Printf("blocklist not loaded: %v", err)
		localBlockDomain = make(map[string]struct{})
	} else {
		log.Printf("blocklist loaded: %d domains", len(localBlockDomain))
	}

	// Spwan workers
	jobs := make(chan jobRequest, 50)
	for range 4 {
		go dnsWorker(localDnsCache, localBlockDomain, conn, jobs)
	}

	log.Println("Server up and running")

	for {

		// listen for the DNS request on the wire. This is blocking
		rawInput := make([]byte, 512)
		rawInputSize, clientAddr, err := conn.ReadFromUDP(rawInput)
		if err != nil {
			log.Printf("Failed to read UDP packet: %v", err)
			continue
		}

		data := make([]byte, rawInputSize)
		copy(data, rawInput[:rawInputSize])

		jobs <- jobRequest{clientAddr, data, rawInputSize}

	}

}

func dnsWorker(localDnsCache map[string][]DNSLocalCache, localBlockDomain map[string]struct{}, conn *net.UDPConn, jobs chan jobRequest) {

	for {

		var (
			dnsRequest, dnsResponse = DNS{}, DNS{}
			bytesToSend             []byte
		)

		// Listen from job Channel
		newJob := <-jobs

		// Convert input bytes to dnsRequest of type DNS
		inputBuff := bytes.NewBuffer(newJob.data[:newJob.size])

		dnsRequest, err := ExtractRequest(inputBuff)
		if err != nil {
			log.Printf("Unable to extract request. Bytes: %v", inputBuff)
			continue
		}

		// discard request if cetain conditions are met
		if DiscardRequest(&dnsRequest) {
			continue
		}

		// Fetch record from local cache
		recordStat, recordFromLocalCache := FetchRecordFromLocal(localDnsCache, localBlockDomain, &dnsRequest)

		if recordStat == DOMAIN_NOT_FOUND_IN_LOCAL {
			// Send the raw request to Public DNS
			bytesToSend, err = FetchFromPublicDNS(newJob.data, newJob.size)
			if err != nil {
				log.Printf("Failed to fetch from public DNS: %v", err)
				continue
			}
		} else {
			// Craft the response if domain found locally
			SetHeadersAndFields(&dnsRequest, &dnsResponse, recordStat)
			CraftResponseRecord(&dnsResponse, recordStat, recordFromLocalCache)
			bytesToSend = dnsResponse.ToBytes(recordStat)
		}

		_, err = conn.WriteToUDP(bytesToSend, newJob.addr)
		if err != nil {
			log.Printf("Failed to write UDP response: %v", err)
		}

	}

}
