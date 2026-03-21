package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func LoadAllDNSCache() map[string][]DNSLocalCache {

	db := make(map[string][]DNSLocalCache)

	data, err := os.ReadFile("./records.json")
	if err != nil {
		log.Println("no local records found")
		return nil
	}

	err = json.Unmarshal(data, &db)
	if err != nil {

		log.Println(errors.New("invalid record format, skipping records.json"))
		return nil
	}

	log.Println("local records loaded")

	return db

}

func packetBinaryWrite(buff io.Writer, data ...any) {
	// writing data to bytes
	for _, iData := range data {
		binary.Write(buff, binary.BigEndian, iData)
	}
}

func FetchFromPublicDNS(inputBuffer []byte, inputBuffSize int) ([]byte, error) {

	dnsServerBuffer := make([]byte, 512)

	serverAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		return dnsServerBuffer, err
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return dnsServerBuffer, err
	}
	defer conn.Close()

	_, err = conn.Write(inputBuffer[:inputBuffSize])
	if err != nil {
		return dnsServerBuffer, err
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := conn.ReadFromUDP(dnsServerBuffer)
	if err != nil {
		return dnsServerBuffer, err
	}

	return dnsServerBuffer[:n], err

}
