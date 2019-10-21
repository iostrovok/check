/*
	For more details see https://github.com/2tvenom/go-test-teamcity
*/

package check

import (
	"fmt"
	"strings"
	"time"
)

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

const (
	TeamcityTimestampFormat = "2006-01-02T15:04:05.000"
)

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

func teamcityOutput(status string, test *C, details ...string) string {
	now := timeFormat(time.Now())
	testName := escape(*formatMessageNamePrefixFlag + test.testName)

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
