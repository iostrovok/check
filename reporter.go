package check

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

const (
	TeamcityTimestampFormat = "2006-01-02T15:04:05.000"
)

type outputWriter struct {
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	Stream               bool
	Verbose              bool
}

func newOutputWriter(writer io.Writer, stream, verbose bool) *outputWriter {
	return &outputWriter{writer: writer, Stream: stream, Verbose: verbose}
}

func (ow *outputWriter) Write(content []byte) (n int, err error) {
	ow.m.Lock()
	n, err = ow.writer.Write(content)
	ow.m.Unlock()
	return
}

func (ow *outputWriter) WriteCallStarted(label string, c *C) {
	if ow.Stream {
		header := renderCallHeader(label, c, "", "\n")
		ow.m.Lock()
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	}
}

func (ow *outputWriter) WriteCallProblem(label string, c *C) {
	var prefix string
	if !ow.Stream {
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	}
	header := renderCallHeader(label, c, prefix, "\n\n")
	ow.m.Lock()
	ow.wroteCallProblemLast = true
	ow.writer.Write([]byte(header))
	if !ow.Stream {
		c.logb.WriteTo(ow.writer)
	}
	ow.m.Unlock()
}

func (ow *outputWriter) WriteCallSuccess(label string, c *C) {
	if ow.Stream || (ow.Verbose && c.kind == testKd) {
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" {
			suffix = " (" + c.reason + ")"
		}
		if c.status() == succeededSt {
			suffix += "\t" + c.timerString()
		}
		suffix += "\n"
		if ow.Stream {
			suffix += "\n"
		}
		header := renderCallHeader(label, c, "", suffix)
		ow.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if !ow.Stream && ow.wroteCallProblemLast {
			header = "\n-----------------------------------" +
				"-----------------------------------\n" +
				header
		}
		ow.wroteCallProblemLast = false
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	}
}

func renderCallHeader(label string, c *C, prefix, suffix string) string {
	pc := c.method.PC()

	out := fmt.Sprintf("%s%s %s: %s%s", prefix, label, niceFuncPath(pc),
		niceFuncName(pc), suffix)

	if *newMessageFlag {
		out += outputTest(label, c, niceFuncPath(pc), niceFuncName(pc), suffix) +"\n"
	}

	return out

}

func timeFormat(t time.Time) string {
	return t.Format(TeamcityTimestampFormat)
}

func escapeLines(lines []string) string {
	return escape(strings.Join(lines, "\n"))
}

func escape(s string) string {
	s = strings.Replace(s, "|", "||", -1)
	s = strings.Replace(s, "\n", "|n", -1)
	s = strings.Replace(s, "\r", "|n", -1)
	s = strings.Replace(s, "'", "|'", -1)
	s = strings.Replace(s, "]", "|]", -1)
	s = strings.Replace(s, "[", "|[", -1)
	return s
}

func outputTest(status string, test *C, details ...string) string {
	now := timeFormat(time.Now())
	testName := escape(*addTestName + test.testName)

	if status == "START" {
		return fmt.Sprintf("##teamcity[testStarted timestamp='%s' name='%s' captureStandardOutput='true']", timeFormat(test.startTime), testName)
	}

	if status == "SKIP" || status == "MISS" {
		return fmt.Sprintf("##teamcity[testIgnored timestamp='%s' name='%s']", now, testName)

	}

	out := ""
	switch status {
	//case "Race":
	//	out += fmt.Sprintf("##teamcity[testFailed timestamp='%s' name='%s' message='Race detected!' details='%s']\n",
	//		now, testName, escapeLines(details))
	case "FAIL":
		out += fmt.Sprintf("##teamcity[testFailed timestamp='%s' name='%s' details='%s']",
			now, testName, escapeLines(details))
	case "PASS":
		// ignore
	default:
		out += fmt.Sprintf("##teamcity[testFailed timestamp='%s' name='%s' message='Test ended in panic.' details='%s']",
			now, testName, escapeLines(details))
	}

	out += fmt.Sprintf("##teamcity[testFinished timestamp='%s' name='%s' duration='%d']",
		now, testName, test.duration/time.Millisecond)

	return out
}
