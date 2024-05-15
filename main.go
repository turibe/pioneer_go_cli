package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
		text, e := reader.ReadString('\n')
		if e == io.EOF {
			exit()
		}
		if len(text) > 1 {
			// send to socket
			command := strings.ToLower(strings.TrimSpace(text))
			if command == "quit" || command == "exit" {
				exit()
			}
			if command == "status" {
				get_status(ch)
				continue
			}
			if command == "video status" {
				ch <- "?VST"
				continue
			}
			i, err := strconv.Atoi(command)
			if err == nil {
				if i > 0 {
					fmt.Printf("Volume up %d\n", i)
					i = min(i, 10)
					for x := 0; x < i; x++ {
						ch <- "VU"
					}
				}
				if i < 0 {
					i = Abs(max(i, -30))
					fmt.Printf("Volume down %d\n", i)
					for x := 0; x < i; x++ {
						ch <- "VD"
					}
				}
				continue
			}
			comm := commandMap[command]
			if comm != "" {
				if DEBUG {
					fmt.Printf("Mapped to %s\n", comm)
				}
				text = comm
				// fmt.Fprintf(conn, text+"\n")
			} else {
				fmt.Printf("Unknown command, sending raw: %s\n", command)
				text = command
			}
			ch <- text
		}
	}
}

func exit() {
	fmt.Println("Adios!")
	os.Exit(0)
}

func sender(conn net.Conn, c chan string) {
	for {
		s := <-c
		// fmt.Printf("Got message %s\n", s)
		fmt.Fprintf(conn, s+"\r\n")
	}
}

// sends query commands to get various system status info back
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
			if len(f) > 0 {
				fmt.Println(f)
			}
		} else {
			fmt.Printf("Could not decode %s\n", message)
			fmt.Printf("Message from server, len %d: %s ", len(message), message)
		}
	}
}

// handles the message that comes back from the AVR
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
	if strings.HasPrefix(message, "VOL") {
		return "", nil
	}
	return message, nil
}

// ======= helpers

func Abs[T ~int | ~int32 | ~int64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
