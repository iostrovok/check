// Package check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//     http://labix.org/gocheck
//
package check

// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

var StdLabels map[string]string

func init() {
	StdLabels = map[string]string{
		"FAIL EXPECTED": "--- FAIL: ",
		"FAIL":          "--- FAIL: ",
		"MISS":          "--- SKIP: ",
		"PANIC":         "--- FAIL: ",
		"PASS":          "--- PASS: ",
		"SKIP":          "--- SKIP: ",
		"START":         "=== RUN",
	}
}

func styleLabel(label string) string {
	if *newMessageFlag {
		if m, find := StdLabels[label]; find {
			return m
		}
	}

	return label + ":"
}
