package agent

import "regexp"

// ansiRegex matches ANSI escape sequences: CSI, OSC, charset, and other common codes.
var ansiRegex = regexp.MustCompile(`\x1b(?:\[[0-9;?]*[a-zA-Z@]|\][^\x07\x1b]*(?:\x07|\x1b\\)|[()][AB012]|[>=]|[78DEHM])|\x07`)

// stripAnsi removes all ANSI escape sequences from a string.
func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
