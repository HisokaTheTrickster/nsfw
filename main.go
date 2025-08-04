package main

import (
	"bytes"
	"fmt"
	"net"
	"nsfw/utils"
)

const (
	ADDRESS_PORT = ":53"
)

func raisePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func raiseError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", ADDRESS_PORT)
	raisePanic(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	raisePanic(err)
	defer conn.Close()

	fmt.Println("Starting DNS server ...")

	// buffer to recieve the message
	inputBuff := make([]byte, 512)
	//outputBuff := make([]byte, 512)

	for noOfRequest := 0; noOfRequest < 1; noOfRequest++ {

		fmt.Println("Waiting for Requests ...")

		n, clientAddr, err := conn.ReadFromUDP(inputBuff)
		raisePanic(err)

		fmt.Printf("Sending Request to Input Handler Length of Packet %d, DNS Requested by: %v\n", n, clientAddr)
		err = utils.DNSRequestHandler(bytes.NewBuffer(inputBuff[:n]))
		raiseError(err)

		// fmt.Errorf("Send Response to cl ient")
		// DNSResponseHandler(clientAddr, outputBuff)

	}

}
