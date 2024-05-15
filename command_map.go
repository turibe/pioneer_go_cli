package main

var commandMap = map[string]string{
	"on":     "PO",
	"off":    "PF",
	"up":     "VU",
	"+":      "VU",
	"down":   "VD",
	"-":      "VD",
	"mute":   "MO",
	"unmute": "MF",

	"volume": "?V",

	"tone":         "9TO", // cyclic
	"tone off":     "0TO",
	"tone on":      "1TO",
	"treble up":    "TI",
	"treble down":  "TD",
	"treble reset": "06TR",
	"bass up":      "BI",
	"bass down":    "BD",
	"bass reset":   "06BA",

	"mcacc": "MC0", // cyclic

	// phase control is recommended to be on:
	"phase": "IS9", // cyclic

	"stereo":    "0001SR", // cycle through stereo modes
	"unplugged": "0109SR",
	"extended":  "0112SR",

	"mode":             "?S",
	"model":            "?RGD",
	"mac address":      "?SVB",
	"software version": "?SSI",

	"loud": "9ATW", // cyclic

	// commands for switching inputs have the form XXFN, now derived from inputSourcesMap.
	// "phono" : "00FN", # invalid command
	// "hdmi" : "31FN", # cyclic

	// TODO: could have a pandora mode, radio mode, etc.
	// Pandora ones:
	"start":    "30NW",
	"next":     "13NW",
	"pause":    "11NW",
	"play":     "10NW",
	"previous": "12NW",
	"stop":     "20NW",
	"clear":    "33NW",
	"repeat":   "34NW",
	"random":   "35NW",
	"menu":     "36NW",

	"info":     "?GAH",
	"list":     "?GAI",
	"top menu": "19IP",

	// Tuner ones:
	"nextpreset": "TPI",
	"prevpreset": "TPD",
	"mpx":        "05TN",

	// Cyclic mode shortcuts:
	// cycles through thx modes, but input must be THX:
	"thx": "0050SR",
	// cycles through surround modes (shortcut for "mode" command):
	"surr":         "0100SR",
	"video status": "?VST",
}
