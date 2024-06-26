package main

import (
	"bufio"
	"cmp"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	// _ "net/http/pprof"
)

var DEBUG = false

var print_mutex sync.Mutex

//go:embed commandMap.json
var MainFolder embed.FS

func report(format string, args ...interface{}) (int, error) {
	print_mutex.Lock()
	defer print_mutex.Unlock()
	return fmt.Printf(format, args...)
}

func readCommandMap() map[string]([]string) {
	var err error
	var data []byte
	filename := "commandMap.json"
	data, err = MainFolder.ReadFile(filename)
	if err != nil {
		report("Could not read command map from %s: %s\n", filename, err)
		os.Exit(1)
	}
	var cmap map[string]([]string)
	err = json.Unmarshal(data, &cmap)
	if err != nil {
		report("Error parsing json command map from %s: %v\n", filename, err)
		os.Exit(1)
	}
	report("Read command map from %s\n", filename)
	return cmap
}

var commandMap map[string]([]string)

func init() {
	commandMap = readCommandMap()
}

func main() {

	/*** for benchmarking:
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	***/

	args := os.Args
	var address string
	if len(args) > 1 {
		// report("Found args %v\n", args)
		address = args[1]
	} else {
		address = "192.168.86.32:23"
	}
	regex, ec := regexp.Compile("(.*):([0-9]+)")
	if ec != nil {
		report("error compiling regex: %s\n", ec)
		os.Exit(1)
	}
	matches := regex.FindAllStringSubmatch(address, -1)
	default_port := "23"
	if len(matches) == 0 {
		report("Using default port %s\n", default_port)
		address = strings.Join([]string{address, default_port}, ":")
	} else {
		// report("Found port %s\n", matches)
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		report("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	report("Connected to %s\n", address)

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
		text = strings.TrimSpace(text)
		command := strings.ToLower(text)
		if command == "" {
			continue
		}
		// note: if string is empty, split_command returns [empty]
		split_command := strings.Split(command, " ")
		if DEBUG {
			report("%s %d\n", split_command, len(split_command))
		}
		if len(split_command) == 0 {
			continue
		}
		base_command := split_command[0]
		/*
			second_arg := ""
			if len(split_command) > 1 {
				second_arg = split_command[1]
			}
		*/
		i, err := strconv.Atoi(command)
		pair, ok := commandMap[command]
		var comm string
		if !ok {
			comm = SOURCE_MAP.inverse_map[command]
		} else {
			comm = pair[0]
		}
		switch {
		case command == "quit" || command == "exit":
			exit()
		case command == "debug":
			DEBUG = !DEBUG
			report("Debug is now %v\n", DEBUG)
		case command == "status":
			getStatus(ch)
		case command == "learn":
			for i := 0; i < 60; i++ {
				s := fmt.Sprintf("?RGB%2d", i)
				ch <- s
			}
		case command == "save":
			SOURCE_MAP.save_to_file()
		case command == "sources" || command == "inputs":
			printInputSourceHelp()
		case command == "modes":
			printModeHelp()
		case base_command == "help" || base_command == "?":
			if command == base_command {
				print_help()
				continue
			}
			if len(split_command) == 2 {
				topic := split_command[1]
				if strings.HasPrefix("modes", topic) {
					printModeHelp()
					continue
				}
				if strings.HasPrefix("sources", topic) {
					printInputSourceHelp()
					continue
				}
				pair, ok := commandMap[topic]
				if ok && len(pair) == 2 {
					report("%s: %s\n", topic, pair[1])
					continue
				} else {
					report("Unknown command: %s\n", topic)
					continue
				}
			}
			report("Couldn't recognize help command %s\n", command)
			continue
		// skipping "select" and "display" for now
		case err == nil:
			if i > 0 {
				report("Volume up %d\n", i)
				i = min(i, 10)
				for x := 0; x < i; x++ {
					ch <- "VU"
				}
			}
			if i < 0 {
				i = Abs(max(i, -30))
				report("Volume down %d\n", i)
				for x := 0; x < i; x++ {
					ch <- "VD"
				}
			}

		case comm != "":
			if DEBUG {
				report("Mapped to %s\n", comm)
			}
			for _, text := range strings.Split(comm, ",") {
				ch <- strings.TrimSpace(text)
			}
		case base_command == "mode":
			changeMode(ch, split_command)
		default:
			report("Unknown command, sending raw: %s\n", text)
			ch <- text
		}
	}
}

func changeMode(ch chan string, split_command []string) bool {
	if len(split_command) < 2 {
		return false
	}
	modestring := strings.Join(split_command[1:], " ")
	if modestring == "help" {
		printModeHelp()
		return true
	}
	mset := get_modes_with_prefix(modestring)
	if len(mset) == 1 {
		mode := mset[0]
		m := inverseModeSetMap[mode]
		report("Trying to change mode to %s (%s)\n", modestring, m)
		ch <- m + "SR"
		return true
	}
	print_mutex.Lock()
	defer print_mutex.Unlock()
	fmt.Println("Which mode do you mean? Options are:")
	for i := 0; i < len(mset); i++ {
		println(mset[i])
	}
	return false
}

// Lists the mode change options (not all work)
func printModeHelp() {
	print_mutex.Lock()
	defer print_mutex.Unlock()
	fmt.Println("mode [mode]\tfor one of:")
	for _, k := range SortedKeys(inverseModeSetMap) {
		println(k)
	}
}

// Returns keys that start with "prefix" --- with the special case
// where if "string" is a key, only {string} is returned.
func get_modes_with_prefix(prefix string) []string {
	r := []string{}
	_, ok := inverseModeSetMap[prefix]
	if ok {
		return []string{prefix}
	}
	for k := range inverseModeSetMap {
		if strings.HasPrefix(k, prefix) {
			r = append(r, k)
		}
	}
	return r
}

// prints help for input source change options
func printInputSourceHelp() {
	print_mutex.Lock()
	defer print_mutex.Unlock()
	fmt.Println("Enter one of the following to change input:")
	for _, k := range SortedKeys(SOURCE_MAP.inverse_map) {
		fmt.Printf("%s (%s)\n", k, SOURCE_MAP.inverse_map[k])
	}
	fmt.Print("Use 'learn' to update this map, 'save' to save it to a JSON file.\n")
}

func exit() {
	report("Bye! Adiós!\n")
	os.Exit(0)
}

func sender(conn net.Conn, c <-chan string) {
	var s string
	last := time.Now()
	for {
		s = strings.TrimSpace(<-c)
		if DEBUG {
			fmt.Printf("Got message %s\n", s)
		}
		t := time.Since(last)
		time.Sleep(max(10*time.Millisecond-t, 0))
		fmt.Fprintf(conn, s+"\r\n")
		last = time.Now()
	}
}

// sends query commands to get various system status info back
func getStatus(c chan<- string) {
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
			report("Read error: %v", e)
		}
		if DEBUG {
			report("Read from conn: %s", message)
		}
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}
		f, e := decode_message(message)
		if e == nil && len(f) > 0 {
			report("%s\n", f)
		} else {
			if e != nil {
				report("Could not decode %v\n", e)
				report("Message from server, len %d: %s ", len(message), message)
			}
		}
	}
}

// handles the message that comes back from the AVR
func decode_message(message string) (string, error) {
	if DEBUG {
		report("Decoding %s\n", message)
	}
	// TODO: this function often checks things twice, builds unecessary errors.
	// Ideally would simplify to a switch.
	em := ErrorMap[message]
	if em != "" {
		return em, nil
	}
	if strings.HasPrefix(message, "RGB") {
		SOURCE_MAP.learn_input_from(message[3:])
		return "", nil
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
	f, e = decode_vta(message)
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

	if strings.HasPrefix(message, "SVB") {
		report("AVR mac address: %s\n", message[3:])
		return "", nil
	}
	if strings.HasPrefix(message, "SSI") {
		report("AVR software version: %s\n", message[3:])
		return "", nil
	}

	if strings.HasPrefix(message, "FN") {
		// fmt.Printf("Got input %s\n", message)
		inputString := SOURCE_MAP.source_map[message[2:]]
		return fmt.Sprintf("Input is %s", inputString), nil
	}
	if strings.HasPrefix(message, "ATW") {
		f, e = SwitchChar(message, 3)
		return fmt.Sprintf("loudness is %s", f), e
	}
	if strings.HasPrefix(message, "ATC") {
		f, e = SwitchChar(message, 3)
		return fmt.Sprintf("eq is %s", f), e
	}
	if strings.HasPrefix(message, "ATD") {
		f, e = SwitchChar(message, 3)
		return fmt.Sprintf("standing wave is %s", f), e
	}
	if strings.HasPrefix(message, "ATE") {
		num := message[3:]
		switch {
		case "00" <= num && num <= "16":
			return fmt.Sprintf("Phase control: %s ms", num), nil
		case num == "97":
			return "Phase control: AUTO", nil
		case num == "98":
			return "Phase control: UP", nil
		case num == "99":
			return "Phase control: DOWN", nil
		default:
			return "Phase control: unknown", nil
		}
	}
	f, e = translate_mode(message)
	if e == nil {
		return f, e
	}
	if strings.HasPrefix(message, "AST") {
		return decode_ast(message)
	}
	if strings.HasPrefix(message, "VTC") {
		return decode_vtc(message)
	}
	if strings.HasPrefix(message, "SR") {
		code := message[2:]
		v, ok := modeSetMap[code]
		if ok {
			return fmt.Sprintf("mode is %s (%s)", v, message), nil
		}
	}
	f, e = decode_vst(message)
	if e == nil {
		return f, e
	}
	if strings.HasPrefix(message, "RGD") {
		return fmt.Sprintf("AVR model info: %s", message), nil
	}
	if strings.HasPrefix(message, "VOL") {
		vol := volDbLevel(message[3:])
		return fmt.Sprintf("volume is %s", vol), nil
	}
	return message, fmt.Errorf("unknown message %s", message)
}

func print_help() {
	commands := []string{"help", "status", "quit", "learn", "debug", "save"}
	k := Keys(commandMap)
	commands = slices.Concat(commands, k)
	slices.Sort(commands)
	for _, c := range commands {
		println(c)
	}
	print("Use \"help mode\" for information on modes, \"help sources\" for changing input sources, \"quit\" to exit\n")
}

// ======= helper functions:

func Abs[T ~int | ~int32 | ~int64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Keys[K comparable, T any](m map[K]T) (result []K) {
	result = make([]K, len(m))
	count := 0
	for k := range m {
		result[count] = k
		count += 1
	}
	return result
}
func SortedKeys[K cmp.Ordered, T any](m map[K]T) (result []K) {
	result = Keys(m)
	slices.Sort(result)
	return result
}

func SwitchChar(s string, i int) (string, error) {
	if i >= len(s) {
		return "", fmt.Errorf("switch error, string too short: %s, %d", s, i)
	}
	switch s[i] {
	case '1':
		return "on", nil
	case '0':
		return "off", nil
	default:
		return "", fmt.Errorf("unexpected switch: %s, %d", s, i)
	}
}
