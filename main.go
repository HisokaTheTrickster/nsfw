package main

import (
	"bytes"
	"log"
	"net"
)

func raisePanic(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func main() {

	serverAddr, _ := net.ResolveUDPAddr("udp", DNS_ADDRESS_PORT)
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	log.Printf(`



===============================================================

	 /$$   /$$  /$$$$$$  /$$$$$$$$ /$$      /$$
	| $$$ | $$ /$$__  $$| $$_____/| $$  /$ | $$
	| $$$$| $$| $$  \__/| $$      | $$ /$$$| $$
	| $$ $$ $$|  $$$$$$ | $$$$$   | $$/$$ $$ $$
	| $$  $$$$ \____  $$| $$__/   | $$$$_  $$$$
	| $$\  $$$ /$$  \ $$| $$      | $$$/ \  $$$
	| $$ \  $$|  $$$$$$/| $$      | $$/   \  $$
	|__/  \__/ \______/ |__/      |__/     \__/


       		NSFW - Name Server For the Web
  A lightweight, fast, and customizable DNS resolver for lan

   	  https://github.com/your-username/nsfw	
                                             
===============================================================
                                           
	`)

	log.Println("loading local records")
	localDnsCache := LoadAllDNSCache()

	// buffer to recieve the message


	log.Println("Server up and running")

	for {

		var (
			dnsRequest, dnsResponse = DNS{}, DNS{}
			bytesToSend             []byte
		)

		// listen for the DNS request on the wire. This is blocking
		rawInput := make([]byte, 512)
		rawInputSize, clientAddr, _ := conn.ReadFromUDP(rawInput)
		

		// Convert input bytes to dnsRequest of type DNS
		inputBuff := bytes.NewBuffer(rawInput[:rawInputSize])
		dnsRequest, err := ExtractRequest(inputBuff)
		if err != nil {
			continue
		}

		// discard request if cetain conditions are met
		if DiscardRequest(&dnsRequest) {
			continue
		}

		// Fetch record from local cache
		recordStat, recordFromLocalCache := FetchRecordFromLocal(localDnsCache, &dnsRequest) 

		if recordStat == DOMAIN_NOT_FOUND_IN_LOCAL {
			// Send the raw request to Google DNS
			bytesToSend, err = FetchFromNet(&rawInput, rawInputSize)
		} else {
			SetHeadersAndFields(&dnsRequest, &dnsResponse, recordStat)
			CraftResponseRecord(&dnsResponse, recordStat, recordFromLocalCache )
			bytesToSend = dnsResponse.ToBytes(recordStat)
		}

		_, err = conn.WriteToUDP(bytesToSend, clientAddr)

		raisePanic(err)

	}

}
