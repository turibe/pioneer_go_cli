package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

/****
func main() {

	var caller telnet.Caller = telnet.StandardCaller

	err := telnet.DialToAndCall("192.168.86.32:23", caller)
	if nil != err {
		panic(err)
	}
}
****/

func main() {
	// connect to this socket
	address := "192.168.86.32:23"
	conn, _ := net.Dial("tcp", address)
	fmt.Printf("Connected to %s\n", address)

	defer conn.Close()

	go read(conn)

	ch := make(chan string)

	go sender(conn, ch)

	for {
		// read from stdin:
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
			comm := commandMap[trimmed]
			if comm != "" {
				fmt.Printf("Mapped to %s\n", comm)
				text = comm + "\r\n"
				// fmt.Fprintf(conn, text+"\n")
			} else {
				fmt.Printf("Unknown command %s, seding raw\n", trimmed)
				text = trimmed
			}
			ch <- text
		}
	}
}

func sender(conn net.Conn, c chan string) {
	for {
		s := <-c
		fmt.Fprintf(conn, s+"\n")
	}
}

// listen for reply
func read(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if len(message) == 0 {
			continue
		}
		f, e := decode_fl(message)
		if e == nil {
			fmt.Println(f)
		} else {
			fmt.Printf("Could not decode %s\n", message)
			fmt.Printf("Message from server, len %d: %s ", len(message), message)
		}
	}
}

func decode_fl(s string) (string, error) {
	fmt.Printf("Decoding %s\n", s)
	if !(strings.HasPrefix(s, "FL")) {
		return "", fmt.Errorf("string does not start with FL: %s", s)
	}
	s = s[2:]
	s = s[2:]
	b, e := hex.DecodeString(s)
	return string(b), e
}
