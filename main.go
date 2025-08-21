package main

import (
	"bytes"
	"log"
	"net"
	"nsfw/utils"
	"os"
)

const (
	ADDRESS_PORT = ":53"
)

func raisePanic(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func main() {

	logFile, err := os.OpenFile("logs/dns.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Setting up the Server ...")

	udpAddr, err := net.ResolveUDPAddr("udp", ADDRESS_PORT)
	raisePanic(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	raisePanic(err)
	defer conn.Close()

	allDnsCache, err := utils.LoadAllDNSCache()
	raisePanic(err)

	// buffer to recieve the message
	inputBuff := make([]byte, 512)
	//outputBuff := make([]byte, 512)

	for {

		log.Println("Waiting for Requests ...")

		n, clientAddr, _ := conn.ReadFromUDP(inputBuff)
		log.Printf("Request recived\nExtracting headers... \nLength of Packet %d, DNS Requested by: %v\n", n, clientAddr)

		dnsRequest, err := utils.ExtractRequest(bytes.NewBuffer(inputBuff[:n]))
		if err != nil {
			log.Println(err.Error() + ". dropping the packet")
			continue
		}

		// discard packets if cetain conditions are met
		if utils.DiscardRequest(&dnsRequest) {
			continue
		}

		bytesToSend := bytes.Buffer{}
		dnsResponse := utils.DNS{}
		err = utils.FetchFromLocalRecord(allDnsCache, &dnsRequest, &dnsResponse)

		if err == nil {
			utils.CopyRequiredFields(&dnsRequest, &dnsResponse)
			bytesToSend = dnsResponse.ToBytes()
		} else {
			if err == utils.ErrNoLocalRecord {
				// send packet to Google DNS
				log.Println("Need to send it to Googld DNS")
				continue

			} else {
				log.Println(err)
				continue
			}
		}

		log.Println("Sending response")
		_, err = conn.WriteToUDP(bytesToSend.Bytes(), clientAddr)

		raisePanic(err)

	}

}
