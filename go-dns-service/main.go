package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/miekg/dns"
)

const (
	DNS_PORT = 53
)

type DbConfig struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Passwd   string `json:"passwd"`
	Db   	 int    `json:"db"`
}

func initDbClient() (*redis.Client, error) {
	// Load the Redis configuration from a JSON file.
	configBytes, err := ioutil.ReadFile("db_config.json")
	if err != nil {
		panic(err)
	}
	var config DbConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	// Start Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     config.Hostname + ":" + config.Port,
		Password: config.Passwd, // no password set
		DB:       config.Db,  // use default DB
	})

	return client, err
}

func main() {
	// Connect to Redis server
	dbClient, err := initDbClient()

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
		go handleRequest(serverConn, addr, buf[:n], dbClient)
	}
}

func delegate(conn *net.UDPConn, addr *net.UDPAddr, request []byte) error {
	// Delegate request to 8.8.8.8
	delegationServerAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = conn.WriteToUDP(request, delegationServerAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Wait for response from 8.8.8.8
	response := make([]byte, 65535)
	n, _, err := conn.ReadFromUDP(response)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Truncate response to actual length
	response = response[:n]
	fmt.Println("Received response from 8.8.8.8: ", response)

	// Send response back to original client
	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getValueFromRedis(dbClient *redis.Client, key string) (string, error) {
	// create context
	ctx := context.Background()
	// Read from Redis database
	val, err := dbClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val,nil
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte, dbClient *redis.Client) {
	// Parse DNS request
	msg := &dns.Msg{}
	err := msg.Unpack(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// lab domain name
	labDomain := ".lab.local."

	// Domain name being requested
	domainName := msg.Question[0].Name

	// Print domain name being requested
	fmt.Println("Received request for domain:", domainName)

	// Check if the domain name ends with ".lab.local."
	if !strings.HasSuffix(domainName, labDomain) {
		fmt.Println("Delegating request to 8.8.8.8...")
		// Delegate request to 8.8.8.8
		err = delegate(conn, addr, request)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Look up domain name in DNS database
	hostname := domainName[:len(domainName)-len(labDomain)]
	fmt.Println("Looking up IP address for hostname:", hostname)
	ip, err := getValueFromRedis(dbClient, hostname)
	fmt.Println("IP address for hostname:", ip)
	if err != nil {
		fmt.Println(err)
		return
	}


	// Encode DNS response
	// response, err := msg.Pack()
	// response is ip as[]byte
	response := []byte(ip)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Send DNS response
	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		fmt.Println(err)
	}
}
