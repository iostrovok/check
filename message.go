// Package check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//     http://labix.org/gocheck
//
package check

// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

var StdMessages map[]string

func init() {
	StdMessages = map[]string{
		"FAIL EXPECTED": "--- FAIL",
		"FAIL":          "--- FAIL",
		"MISS":          "--- SKIP",
		"PANIC":         "--- FAIL",
		"PASS":          "--- PASS",
		"SKIP":          "--- SKIP",
		"START":         "=== RUN",
	}
}

func message(oldMessage string) string {
	if *newMessageFlag {
		if m, find := StdMessages[oldMessage]; find {
			return m
		}
	}

	return oldMessage
}
