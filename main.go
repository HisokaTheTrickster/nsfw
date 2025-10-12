package main

import (
	"bytes"
	"log"
	"net"
	"os"

	"github.com/HisokaTheTrickster/nsfw/utils"
)

func raisePanic(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func main() {

	logFile, err := os.OpenFile("nsfw.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	serverAddr, _ := net.ResolveUDPAddr("udp", utils.DNS_ADDRESS_PORT)
	conn, _ := net.ListenUDP("udp", serverAddr)
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
	localDnsCache := utils.LoadAllDNSCache()

	// buffer to recieve the message
	inputBuff := make([]byte, 512)

	log.Println("Server up an running")

	for {

		inputBuffSize, clientAddr, _ := conn.ReadFromUDP(inputBuff)

		var (
			dnsRequest, dnsResponse = utils.DNS{}, utils.DNS{}
			bytesToSend             []byte
			recordStat              utils.RecordStatus
		)

		// Convert bytes to dnsRequest of type DNS
		dnsRequest, err := utils.ExtractRequest(bytes.NewBuffer(inputBuff[:inputBuffSize]))
		if err != nil {
			continue
		}

		// discard request if cetain conditions are met
		if utils.DiscardRequest(&dnsRequest) {
			continue
		}

		recordStat, bytesToSend = utils.FetchRecord(localDnsCache, &dnsRequest, &dnsResponse, &inputBuff, inputBuffSize)

		switch recordStat {

		case utils.RECORD_FOUND_LOCALLY:
			utils.CopyRequiredFields(&dnsRequest, &dnsResponse)
			bytesToSend = dnsResponse.ToBytes()

		case utils.DOMAIN_EXISTS_NO_RECORD:
			utils.CopyRequiredFields(&dnsRequest, &dnsResponse)
			bytesToSend = dnsResponse.ToBytes()

		case utils.RECORD_FROM_REMOTE:

		case utils.ERR_REMOTE_DNS_TIMEOUT:
			// Discard packet
			continue

		default:
			log.Println("Error occured when fetching record")
		}

		_, err = conn.WriteToUDP(bytesToSend, clientAddr)

		raisePanic(err)

	}

}
