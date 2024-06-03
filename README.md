# pioneer_go_cli
Golang version of Pioneer CLI telnet controller

==================

Barebones command-line interface (CLI) for controlling an Internet-connected Pioneer AVR (Receiver/Amp).

Tested on a Pioneer SC-1222-K Amp.

License: MIT.

Disclaimer: *Use at your own risk.*

For a more complete API, see (https://github.com/crowbarz/aiopioneer).

This is a (somewhat useful) toy program to try new languages and language features.
See also [a new Rust implementation of the same functionality here](https://github.com/turibe/pioneer_rust_cli),
and the [Python 3 version](https://github.com/turibe/pioneer_python_cli).

## Usage:

0. `go build` or `go run .`.
1. Find out your AVR's IP address.
2. Run with `<ipaddress>` as the argument.

## Some commands:

- `up`              [volume up]
- `down`            [volume down]
- `<integer>`       [if positive, increase volume this number of times, capped at 10]
- `-<integer>`      [if negative, decrease volume this number of times, capped at -30]

- `<input_name>`    [switch to given input]

- `mode X`          [choose audio modes; not all modes will be available]
- `mode help`       [help with modes]
- `help` or `help <command>`
- `surr`            [cycle through surround modes]
- `stereo`          [stereo mode]
- `status`          [print status]

- Use control-D to exit.

If you have customized the input names for your AVR (for example, "AppleTV", "My DVR", etc.),
`learn` gets them from the AVR, after which they are available as commands.
The `save` command saves a JSON file with these names, to be loaded in future sessions, at startup time.

