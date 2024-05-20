package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

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

var sources_map_filename = "pioneer_avr_sources.json"

// TODO: replicate sources.py

type SourceMap struct {
	source_map  map[string]string
	inverse_map map[string]string
	alias_map   map[string]string
}

var SOURCE_MAP SourceMap

func init() {
	SOURCE_MAP = SourceMap{}
	SOURCE_MAP.inverse_map = map[string]string{}
	SOURCE_MAP.alias_map = map[string]string{}
	SOURCE_MAP.init_from_map(defaultInputSourcesMap)
	SOURCE_MAP.read_from_file()
	SOURCE_MAP.add_aliases()
}

func (m *SourceMap) init_from_map(initmap map[string]string) {
	m.source_map = map[string]string{}
	for k, v := range initmap {
		m.source_map[k] = v
		m.register_reverse_source(k, v)
	}
}

func (m *SourceMap) read_from_file() {
	// TODO: prefer user's home directory first
	dirname, err := os.UserHomeDir()
	if err != nil {
		report("Could not get home directory when looking for custom sources %v", err)
		return
	}
	filename := filepath.Join(dirname, sources_map_filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		report("Could not read json file %s: %v\n", filename, err)
		return
	}
	// var mystruct []interface{}
	var mystruct map[string]string
	err = json.Unmarshal(data, &mystruct)
	if err != nil {
		report("Error parsing json sourcemap from %s: %v\n", filename, err)
		return
	}
	m.init_from_map(mystruct)
	report("Updated sources map from %s\n", filename)
}

func (m *SourceMap) save_to_file() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		report("Couldn't get home directory for saving custom source map: %v", err)
		return
	}
	filename := filepath.Join(dirname, sources_map_filename)
	data, err := json.Marshal(m.source_map)
	if err != nil {
		report("Error building json data: %v\n", err)
	} else {
		err = os.WriteFile(filename, data, 0666)
		if err != nil {
			report("Error writing json file to %s: %s\n", filename, err)
		} else {
			report("Wrote sources map to %s\n", filename)
		}
	}
}

func (m *SourceMap) register_reverse_source(k string, v string) {
	newk := strings.ToLower(v)
	m.inverse_map[newk] = strings.Join([]string{k, "FN"}, "")
}

func (m *SourceMap) update_source(name string, id string) {
	report("Updating source %s (%s)\n", name, id)
	m.source_map[id] = name
	m.register_reverse_source(id, name)
	alias, ok := m.alias_map[name]
	if ok {
		m.check_aliases(name, alias)
	}
}

func (m *SourceMap) add_alias(a string, b string) {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	m.alias_map[a] = b
	m.alias_map[b] = a
	m.check_aliases(a, b)
}
func (m *SourceMap) check_aliases(a string, b string) {
	_, has_a := m.inverse_map[a]
	_, has_b := m.inverse_map[b]
	if !has_a && has_b {
		m.inverse_map[a] = m.inverse_map[b]
	} else {
		if !has_b && has_a {
			m.inverse_map[b] = m.inverse_map[a]
		}
	}
}

func (m *SourceMap) add_aliases() {
	m.add_alias("apple", "appletv")
	m.add_alias("amazon", "amazontv")
	m.add_alias("radio", "tuner")
	m.add_alias("iradio", "internet radio")
}

func (m *SourceMap) learn_input_from(s string) {
	id := s[0:2]
	name := s[3:]
	if m.source_map[id] != name {
		report("Updating source name %s for %s\n", name, id)
		m.update_source(name, id)
	}
}
