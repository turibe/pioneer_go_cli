package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/reiver/go-telnet"
)

func main1() {
	var caller telnet.Caller = telnet.StandardCaller

	err := telnet.DialToAndCall("192.168.86.32:23", caller)
	if nil != err {
		panic(err)
	}
}

func main() {
	// connect to this socket
	address := "192.168.86.32:23"
	conn, _ := net.Dial("tcp", address)
	fmt.Printf("Connected to %s\n", address)

	defer conn.Close()

	go read(conn)

	for {
		// read in input
		// from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		if len(text) > 1 {
			// send to socket
			trimmed := text[0 : len(text)-1]
			if trimmed == "quit" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
			text = trimmed + "\r\n"

			fmt.Fprintf(conn, text+"\n")
		}
	}
}

// listen for reply
func read(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if len(message) > 1 {
			fmt.Print("Message from server: " + message)
		} else {
			fmt.Print("got empty message")
		}
	}
}
