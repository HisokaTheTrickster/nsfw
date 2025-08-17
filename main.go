package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"nsfw/utils"
	"os"
)

const (
	ADDRESS_PORT = ":53"
	DNS_DB_PATH  = "records.json"
)

func raisePanic(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func raiseError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func loadAllDNSCache() (map[string]utils.DNSdatabase, error) {

	db := []utils.DNSdatabase{}

	data, err := os.ReadFile(DNS_DB_PATH)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	dnsCache := make(map[string]utils.DNSdatabase)

	for _, record := range db {
		dnsCache[record.Name] = record
	}

	return dnsCache, nil

}

func setupLogger() error {
	file, err := os.OpenFile("logs/dns.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	log.SetOutput(file)
	return nil
}

func main() {

	err := setupLogger()
	raisePanic(err)

	log.Println("Setting up the Server ...")

	udpAddr, err := net.ResolveUDPAddr("udp", ADDRESS_PORT)
	raisePanic(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	raisePanic(err)
	defer conn.Close()

	//dnsCache := Loading Cache
	// DNS cache is a map of URL and Records

	allDnsCache, err := loadAllDNSCache()
	raisePanic(err)

	// buffer to recieve the message
	inputBuff := make([]byte, 512)
	//outputBuff := make([]byte, 512)

	for {

		log.Println("Waiting for Requests ...")

		n, clientAddr, err := conn.ReadFromUDP(inputBuff)
		raisePanic(err)

		log.Printf("Request recived\n3Extracting headers... \nLength of Packet %d, DNS Requested by: %v\n", n, clientAddr)
		dnsPacket, err := utils.RequestHandler(bytes.NewBuffer(inputBuff[:n]))
		raisePanic(err)

		// Use pointers for allDNSCache later on?
		log.Println("Checking Database for a response")
		err = utils.FetchRecord(allDnsCache, &dnsPacket)
		raisePanic(err)

		log.Println("Printing Response")
		//err = utils.ConstructReponse(queryRecord, &dnsPacket)
		bytesToSend := dnsPacket.ToBytes()
		fmt.Println(bytesToSend)

		_, err = conn.WriteToUDP(bytesToSend.Bytes(), clientAddr)

		raisePanic(err)

	}

}
