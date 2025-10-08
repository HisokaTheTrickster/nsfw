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

	// Setting up the logger to save it int he file.
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	logFile, err := os.OpenFile("logs/dns.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Setting up the Server ...")

	serverAddr, _ := net.ResolveUDPAddr("udp", utils.DNS_ADDRESS_PORT)
	conn, _ := net.ListenUDP("udp", serverAddr)
	defer conn.Close()

	localDnsCache, err := utils.LoadAllDNSCache()
	raisePanic(err)

	// buffer to recieve the message
	inputBuff := make([]byte, 512)

	for {

		log.Println("Waiting for Requests ...")

		inputBuffSize, clientAddr, _ := conn.ReadFromUDP(inputBuff)
		log.Printf("Request recived\nExtracting headers... \nLength of Packet %d, DNS Requested by: %v\n", inputBuffSize, clientAddr)

		var (
			dnsRequest, dnsResponse = utils.DNS{}, utils.DNS{}
			bytesToSend             []byte
			recordStat              utils.RecordStatus
		)

		// Convert bytes to dnsRequest of type DNS
		dnsRequest, err := utils.ExtractRequest(bytes.NewBuffer(inputBuff[:inputBuffSize]))
		if err != nil {
			log.Println(err.Error() + ". dropping the packet")
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

		case utils.ERR_REMOTE_DNS_TIMEOUT:
			// Discard packet
			continue

		default:
			log.Println("Error occured when fetching record")
		}

		log.Println("Sending response")
		_, err = conn.WriteToUDP(bytesToSend, clientAddr)

		raisePanic(err)

	}

}
