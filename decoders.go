package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func decode_fl(s string) (string, error) {
	// fmt.Printf("Decoding %s\n", s)
	if !(strings.HasPrefix(s, "FL")) {
		return "", fmt.Errorf("string does not start with FL: %s", s)
	}
	bs := s[4:]
	// fmt.Printf("bs is %s\n", bs)
	b, e := hex.DecodeString(bs)
	// fmt.Printf("Decoded to %s, %v", b, e)
	return string(b), e
}

func decode_geh(s string) (string, error) {
	if strings.HasPrefix(s, "GDH") {
		sbytes := s[3:]
		return "items " + sbytes[0:5] + " to " + sbytes[5:10] + " of total " + sbytes[10:], nil
	}
	if strings.HasPrefix(s, "GBH") {
		return "max list number: " + s[2:], nil
	}
	if strings.HasPrefix(s, "GCH") {
		return screenTypeMap.get(s[3:5], "unknown") + " - " + s, nil
	}
	if strings.HasPrefix(s, "GHH") {
		source := s[2:]
		return "source: " + internetSourceMap.get(source, "unknown"), nil
	}
	if !(strings.HasPrefix(s, "GEH")) {
		return "", fmt.Errorf("not GEH")
	}
	s = s[3:]
	// line = s[0:2]
	// focus = s[2]
	tstring := s[3:5]
	typeval := trackFieldsMap.get(tstring, "unknown ({tstring})")
	info := s[5:]
	return typeval + ": " + info, nil
}

var internetSourceMap = MyMap[string, string]{
	"00": "Intenet Radio",
	"01": "Media Server",
	"06": "SiriusXM",
	"07": "Pandora",
	"10": "AirPlay",
	"11": "Digital Media Renderer (DMR)",
}

var trackFieldsMap = MyMap[string, string]{
	"20": "Track",
	"21": "Artist",
	"22": "Album",
	"23": "Time",
	"24": "Genre",
	"25": "Chapter Number",
	"26": "Format",
	"27": "Bitrate",
	"28": "Category",
	"29": "Composer1",
	"30": "Composer2",
	"31": "Buffer",
	"32": "Channel",
}

var screenTypeMap = MyMap[string, string]{
	"00": "Message",
	"01": "List",
	"02": "Playing (Play)",
	"03": "Playing (Pause)",
	"04": "Playing (Fwd)",
	"05": "Playing (Rev)",
	"06": "Playing (Stop)",
	"99": "Invalid",
}

var VTC_resolution_map = MyMap[string, string]{
	"00": "AUTO Resolution",
	"01": "PURE Resolution",
	"02": "Reserved Resolution",
	"03": "R480/576 Resolution",
	"04": "720p Resolution",
	"05": "1080i Resolution",
	"06": "1080p Resolution",
	"07": "1080/24p Resolution",
}

// Decodes a VTC (video resolution) status status string
func decode_vtc(s string) (string, error) {
	// assert s.startswith('VTC')
	if !(strings.HasPrefix(s, "VTC")) {
		return "", fmt.Errorf("does not start with VTC: %s", s)
	}
	s = s[3:]

	return VTC_resolution_map.get(s, "unknown VTC resolution"), nil
}

var CHANNEL_MAP = map[int]string{
	5:  "Left",
	6:  "Center",
	7:  "Right",
	8:  "SL",
	9:  "SR",
	10: "SBL",
	11: "S",
	12: "SBR",
	13: "LFE",
	14: "FHL",
	15: "FHR",
	16: "FWL",
	17: "FWR",
	18: "XL",
	19: "XC",
	20: "XR",
}

type MyMap[K comparable, V any] map[K]V

func (m MyMap[K, V]) get(k K, deflt V) V {
	v, ok := m[k]
	if !ok {
		return deflt
	}
	return v
}

// Decodes an AST return status string.
func decode_ast(st string) (string, error) {
	var r = ""
	// assert s.startswith("AST")
	s := []byte(st[3:])
	r += fmt.Sprintf("Audio input signal: %s\n", decode_ais(s[0:2]))
	r += fmt.Sprintf("Audio input frequency: %s\n", decode_aif(s[2:4]))

	// The manual starts counting at 1, so to fix this off-by-one, we do:
	s = append([]byte("-"), s...)

	// channels...
	r += "Input Channels:\n"

	for i, v := range CHANNEL_MAP {
		if i >= len(s) {
			continue
		}
		if s[i] == '1' {
			r += fmt.Sprintf("%s,\n", v)
		}
	}
	r += "\n"
	r += "Output Channels:\n"
	for i, v := range CHANNEL_MAP {
		idx := i + 21
		if idx >= len(s) {
			continue
		}
		if s[idx] == '1' {
			r += fmt.Sprintf("%s,\n", v)
		}
	}
	return r, nil
}

var aif_map = MyMap[string, string]{
	"00": "32kHz",
	"01": "44.1kHz",
	"02": "48kHz",
	"03": "88.2kHz",
	"04": "96kHz",
	"05": "176.4kHz",
	"06": "192kHz",
	"07": "---",
}

func decode_aif(s []byte) string {
	return aif_map.get(string(s), "unknown")
}

func decode_ais(st []byte) string {
	s := string(st)
	switch {
	case "00" <= s && s <= "02":
		return "ANALOG"
	case s == "03" || s == "04":
		return "PCM"
	case s == "05":
		return "DOLBY DIGITAL"
	case s == "06":
		return "DTS"
	case s == "07":
		return "DTS-ES Matrix"
	case s == "08":
		return "DTS-ES Discrete"
	case s == "09":
		return "DTS 96/24"
	case s == "10":
		return "DTS 96/24 ES Matrix"
	case s == "11":
		return "DTS 96/24 ES Discrete"
	case s == "12":
		return "MPEG-2 AAC"
	case s == "13":
		return "WMA9 Pro"
	case s == "14":
		return "DSD->PCM"
	case s == "15":
		return "HDMI THROUGH"
	case s == "16":
		return "DOLBY DIGITAL PLUS"
	case s == "17":
		return "DOLBY TrueHD"
	case s == "18":
		return "DTS EXPRESS"
	case s == "19":
		return "DTS-HD Master Audio"
	case "20" <= s && s <= "26":
		return "DTS-HD High Resolution"
	case s == "27":
		return "DTS-HD Master Audio"
	default:
		return "unknown"
	}
}

// db level conversion for treble, bass:
func dbLevel(s string) string {
	n, e := strconv.Atoi(s)
	if e != nil {
		return "?"
	}
	db := 6 - n
	return fmt.Sprintf("{%d}dB", db)
}

func volDbLevel(s string) string {
	n, e := strconv.Atoi(s)
	if e != nil {
		return "?"
	}
	var db = float32(n-161) / 2.0
	return fmt.Sprintf("%.2fdb", db)
}

// readable version of the tone status
func decode_tone(s string) (string, error) {
	if strings.HasPrefix(s, "TR") {
		return "treble at " + dbLevel(s[2:4]), nil
	}
	if strings.HasPrefix(s, "BA") {
		return "bass at " + dbLevel(s[2:4]), nil
	}
	if s == "TO0" {
		return "tone off", nil
	}
	if s == "TO1" {
		return "tone on", nil
	}
	return "", fmt.Errorf("unknown tone %s", s)
}

var SIGNAL_MAP = MyMap[byte, string]{
	'0': "---",
	'1': "VIDEO",
	'2': "S-VIDEO",
	'3': "COMPONENT",
	'4': "HDMI",
	'5': "Self OSD/JPEG",
}

var SIGNAL_FORMAT_MAP = MyMap[string, string]{
	"00": "---",
	"01": "480/60i",
	"02": "576/50i",
	"03": "480/60p",
	"04": "576/50p",
	"05": "720/60p",
	"06": "720/50p",
	"07": "1080/60i",
	"08": "1080/50i",
	"09": "1080/60p",
	"10": "1080/50p",
	"11": "1080/24p",
	"12": "4Kx2K/24Hz",
	"13": "4Kx2K/25Hz",
	"14": "4Kx2K/30Hz",
	"15": "4Kx2K/24Hz(SMPTE)",
}

var ASPECT_MAP = MyMap[byte, string]{
	'0': "---",
	'1': "4:3",
	'2': "16:9",
	'3': "14:9",
}

// HDMI ONLY
var COLOR_MAP = MyMap[byte, string]{
	'0': "---",
	'1': "RGB Limit",
	'2': "RGB Full",
	'3': "YcbCr444",
	'4': "YcbCr422",
}

// HDMI ONLY
var FORMAT_BIT_MAP = MyMap[byte, string]{
	'0': "---",
	'1': "24bit (8bit*3)",
	'2': "30bit (10bit*3)",
	'3': "36bit (12bit*3)",
	'4': "48bit (16bit*3)",
}

var COLOR_SPACE_MAP = MyMap[byte, string]{
	'0': "---",
	'1': "Standard",
	'2': "xvYCC601",
	'3': "xvYCC709",
	'4': "sYCC",
	'5': "AdobeYCC601",
	'6': "AdobeRGB",
}

func decode_vst(st string) (string, error) {
	if !(strings.HasPrefix(st, "VST")) {
		return "", fmt.Errorf("bad VST: %s", st)
	}
	result := ""
	s := "-" + st[3:] // for off-by-one
	// fmt.Printf("Decoding %s\n", s)
	signal := SIGNAL_MAP.get(s[1], "Unknown")
	result += fmt.Sprint("Signal: ", signal, "\n")
	sformat := SIGNAL_FORMAT_MAP.get(s[2:4], "Unknown")
	result += fmt.Sprint("Input resolution: ", sformat, "\n")
	aspect := ASPECT_MAP.get(s[4], "Unknown")
	result += fmt.Sprint("Aspect: ", aspect, "\n")
	color := COLOR_MAP.get(s[5], "Unknown")
	result += fmt.Sprint("Input color format: ", color, "\n")

	ibit := FORMAT_BIT_MAP.get(s[6], "Unknown")
	result += fmt.Sprint("Input bit (HDMI only): ", ibit, "\n")
	cspace := COLOR_SPACE_MAP.get(s[7], "Unknown")
	result += fmt.Sprint("Input extend color space (HDMI only): ", cspace, "\n")
	oformat := SIGNAL_FORMAT_MAP.get(s[8:10], "Unknown")
	result += fmt.Sprint("Output resolution: ", oformat, "\n")
	oaspect := ASPECT_MAP.get(s[10], "Unknown")
	result += fmt.Sprint("Output aspect: ", oaspect, "\n")
	ocolor := COLOR_MAP.get(s[11], "Unknown")
	result += fmt.Sprint("Output color format (HDMI only): ", ocolor, "\n")
	obit := FORMAT_BIT_MAP.get(s[12], "Unknown")
	result += fmt.Sprint("Output bit (HDMI only): ", obit, "\n")
	ospace := COLOR_SPACE_MAP.get(s[13], "Unknown")
	result += fmt.Sprint("Output extend color space (HDMI only): ", ospace, "\n")
	mrecommend := SIGNAL_FORMAT_MAP.get(s[14:16], "Unknown")
	result += fmt.Sprint("Monitor recommend resolution information: ", mrecommend, "\n")
	mdcolor := FORMAT_BIT_MAP.get(s[16], "Unknown")
	result += fmt.Sprint("Monitor DeepColor: ", mdcolor, "\n")

	// ... TODO
	return result, nil
}

func decode_vta(s string) (string, error) {
	if strings.HasPrefix(s, "VTA") {
		return "TODO: decode VTA", nil
	}
	return "", fmt.Errorf("bad VTA: %s", s)
}
