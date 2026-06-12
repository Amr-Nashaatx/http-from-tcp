package main

import (
	"fmt"
	"net"

	"github.com/Amr-Nashaatx/http-from-tcp/request"
)

func main() {
	socket, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(fmt.Errorf("cannot open listen on connections %w", err))
	}
	fmt.Println("Socket open on :42069")
	defer socket.Close()
	con, conErr := socket.Accept()
	fmt.Printf("connection accepted")
	if conErr != nil {
		panic(conErr)
	}
	req, err := request.RequestFromReader(con)
	fmt.Println(req.RequestLine)
	fmt.Println(req.Headers)

	socket.Close()
}
