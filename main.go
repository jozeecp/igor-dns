package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

const (
	DNS_PORT = 53
)

func main() {
	// Bind to DNS port
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", DNS_PORT))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer serverConn.Close()

	fmt.Printf("Listening for DNS requests on %s\n", serverAddr)

	// Start listening for DNS requests
	buf := make([]byte, 1024)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle DNS request
		go handleRequest(serverConn, addr, buf[:n])
	}
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	// Parse DNS request
	msg := &dns.Msg{}
	err := msg.Unpack(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print domain name being requested
	fmt.Println("Received request for domain:", msg.Question[0].Name)

	// Look up domain name in DNS database
	// ...

	// Encode DNS response
	response, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send DNS response
	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		fmt.Println(err)
	}
}
