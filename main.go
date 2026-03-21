package main

import (
	"bytes"
	"log"
	"net"
)

type jobRequest struct{
	addr *net.UDPAddr
	data []byte
	size int
}

func main() {


	log.Printf(BANNER)

	// Listen on port 53
	serverAddr, _ := net.ResolveUDPAddr("udp", DNS_ADDRESS_PORT)
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	// Load records from records.json
	log.Println("loading local records")
	localDnsCache := LoadAllDNSCache()

	// Spwan workers
	jobs := make(chan jobRequest, 50)
	for range(4) {
		go dnsWorker(localDnsCache, conn, jobs)
	}

	log.Println("Server up and running")

	for {

		// listen for the DNS request on the wire. This is blocking
		rawInput := make([]byte, 512)
		rawInputSize, clientAddr, _ := conn.ReadFromUDP(rawInput)

		data := make([]byte, rawInputSize)
        copy(data, rawInput[:rawInputSize])

		jobs <- jobRequest{clientAddr,rawInput, rawInputSize}

	}

}

func dnsWorker(localDnsCache map[string][]DNSLocalCache, conn *net.UDPConn, jobs chan jobRequest) {

	for {

		var (
			dnsRequest, dnsResponse = DNS{}, DNS{}
			bytesToSend             []byte
		)

		// Listen from job Channel
		newJob := <- jobs

		// Convert input bytes to dnsRequest of type DNS
		inputBuff := bytes.NewBuffer(newJob.data[:newJob.size])

		dnsRequest, err := ExtractRequest(inputBuff)
		if err != nil {
			log.Println("Unable to extract request. Bytes: %v", inputBuff)
			continue
		}

		// discard request if cetain conditions are met
		if DiscardRequest(&dnsRequest) {
			continue
		}

		// Fetch record from local cache
		recordStat, recordFromLocalCache := FetchRecordFromLocal(localDnsCache, &dnsRequest) 

		if recordStat == DOMAIN_NOT_FOUND_IN_LOCAL {
			// Send the raw request to Public DNS
			bytesToSend, err = FetchFromPublicDNS(&newJob.data, newJob.size)
		} else {
			// Craft the response if domain found locally
			SetHeadersAndFields(&dnsRequest, &dnsResponse, recordStat)
			CraftResponseRecord(&dnsResponse, recordStat, recordFromLocalCache )
			bytesToSend = dnsResponse.ToBytes(recordStat)
		}

		_, err = conn.WriteToUDP(bytesToSend, newJob.addr)
		if err != nil {
			log.Println(err.Error())
			panic(err)
		}
		
	}

}
