package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/reiver/go-telnet"
)

func main() {

	s := "FL022020204150504C45545620202020"
	s1, e := decode_fl(s)
	if e == nil {
		println(s1)
	}
	os.Exit(0)

	main1()

	var caller telnet.Caller = telnet.StandardCaller

	err := telnet.DialToAndCall("192.168.86.32:23", caller)
	if nil != err {
		panic(err)
	}
}

func main1() {
	// connect to this socket
	address := "192.168.86.32:23"
	conn, _ := net.Dial("tcp", address)
	fmt.Printf("Connected to %s\n", address)

	defer conn.Close()

	go read(conn)

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

func decode_fl(s string) (string, error) {
	if len(s) <= 2 {
		return "", fmt.Errorf("fl string too short: %s", s)
	}
	if s[0:2] != "FL" {
		return "", fmt.Errorf("string does not start with FL: %s", s)
	}
	s = s[2:]
	s = s[2:]
	i := 0
	urls := ""
	for i < len(s) {
		urls += "%"
		urls += s[i : i+2]
		i += 2
	}
	r, e := url.Parse(urls)
	fmt.Printf("%v\n", r.Path)
	return r.Path, e
}
