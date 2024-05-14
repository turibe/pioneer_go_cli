package main

var defaultInputSourcesMap = map[string]string{
	"00": "PHONO",
	"01": "CD",
	"02": "TUNER",
	"04": "DVD",
	"05": "TV",
	"06": "SAT/CBL",
	"10": "VIDEO",
	"12": "MULTI CH IN",
	"13": "USB-DAC",
	"15": "DVR/BDR",
	"17": "iPod/USB",
	"19": "HDMI1",
	"20": "HDMI2",
	"21": "HDMI3",
	"22": "HDMI4",
	"23": "HDMI5",
	"24": "HDMI6",
	"25": "BD",
	"26": "NETWORK", // cyclic
	"31": "HDMI",    // cyclic
	"33": "ADAPTER PORT",
	"34": "HDMI7",
	"38": "INTERNET RADIO",
	"40": "SiriusXM",
	"41": "PANDORA",
	"44": "MEDIA SERVER",
	"45": "Favorites",
	"47": "DMR",
	"48": "MHL", // device input, not working on test AVR
}

var ErrorMap = map[string]string{
	"E02": "NOT AVAILABLE NOW",
	"E03": "INVALID COMMAND",
	"E04": "COMMAND ERROR",
	"E06": "PARAMETER ERROR",
	"B00": "BUSY",
}
