package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalln(err)
	}
	conn, connErr := net.DialUDP("udp", nil, remoteAddr)
	if connErr != nil {
		log.Fatalln(connErr)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			log.Fatalln(readErr)
		}

		_, writeErr := conn.Write([]byte(line))
		if writeErr != nil {
			log.Fatalln(writeErr)
		}
	}

}
