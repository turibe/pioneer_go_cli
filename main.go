package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var DEBUG = false

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

	ch := make(chan string, 100)
	go sender(conn, ch)

	for {
		// read from stdin:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if len(text) > 1 {
			// send to socket
			trimmed := text[0 : len(text)-1]
			if trimmed == "quit" || trimmed == "exit" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
			if trimmed == "status" {
				get_status(ch)
				continue
			}
			comm := commandMap[trimmed]
			if comm != "" {
				if DEBUG {
					fmt.Printf("Mapped to %s\n", comm)
				}
				text = comm
				// fmt.Fprintf(conn, text+"\n")
			} else {
				fmt.Printf("Unknown command, sending raw: %s\n", trimmed)
				text = trimmed
			}
			ch <- text
		}
	}
}

func sender(conn net.Conn, c chan string) {
	for {
		s := <-c
		// fmt.Printf("Got message %s\n", s)
		fmt.Fprintf(conn, s+"\r\n")
	}
}

// @TODO check they all get processed?
func get_status(c chan string) {
	vals := []string{
		"?P",
		"?F",
		"?BA",
		"?TR",
		"?TO",
		"?L",
		"?AST",
		"?IS",
		"?VST",
	}

	for _, v := range vals {
		c <- v
	}
}

// listen for reply
func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, e := reader.ReadString('\n')
		if e != nil {
			fmt.Println("Read error", e)
		}
		// fmt.Println("Read from conn: ", message)
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}
		f, e := decode_message(message)
		if e == nil {
			fmt.Println(f)
		} else {
			fmt.Printf("Could not decode %s\n", message)
			fmt.Printf("Message from server, len %d: %s ", len(message), message)
		}
	}
}

func decode_message(message string) (string, error) {

	em := ErrorMap[message]
	if em != "" {
		return em, nil
	}
	if strings.HasPrefix(message, "RGB") {
		return "TODO: learn from: " + message, nil
	}
	f, e := decode_tone(message)
	if e == nil {
		return f, e
	}
	f, e = decode_geh(message)
	if e == nil {
		return f, e
	}
	f, e = decode_fl(message)
	if e == nil {
		return f, e
	}
	if strings.HasPrefix(message, "IS") {
		switch message[2] {
		case '0':
			return "Phase control OFF", nil
		case '1':
			return "Phase control ON", nil
		case '2':
			return "Full band phase control on", nil
		default:
			return ("Phase control: unknown"), nil
		}
	}

	switch message {
	case "PWR0":
		return "Power is ON", nil
	case "PWR1":
		return "Power is OFF", nil
	}

	if strings.HasPrefix(message, "FN") {
		// println("Got input %s", message)
		inputstring := defaultInputSourcesMap[message[2:]]
		return inputstring, nil
	}
	if strings.HasPrefix(message, "VTC") {
		return decode_vtc(message)
	}
	f, e = translate_mode(message)
	if e == nil {
		return f, e
	}
	if strings.HasPrefix(message, "AST") {
		// TODO: clean up
		decode_ast(message)
		return "", nil
	}
	f, e = decode_vst(message)
	if e == nil {
		return f, e
	}
	return message, nil
}
