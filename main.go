package main

import (
	"bytes"
	"fmt"
	"net"
	"nsfw/utils"
)

const DEF_DNS_FLAG uint16 = 0

func main() {

	address := ":53"
	udpAddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Println("Starting DNS server ...")

	// buffer to recieve the message
	inputBuff := make([]byte, 512)
	//outputBuff := make([]byte, 512)

	for i := 0; i < 1; i++ {

		fmt.Println("Waiting for Requests ...")

		n, clientAddr, err := conn.ReadFromUDP(inputBuff)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Sending Request to Input Handler ")
		fmt.Printf("Length of Packet %d, DNS Requested by: %v\n", n, clientAddr)
		//fmt.Println(n, clientAddr)
		err = utils.DNSRequestHandler(bytes.NewBuffer(inputBuff[:n]))

		if err != nil {
			panic(err)
		}

		// fmt.Errorf("Send Response to cl ient")
		// DNSResponseHandler(clientAddr, outputBuff)

	}

}
