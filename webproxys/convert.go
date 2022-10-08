package webproxys

import (
	"regexp"
	"strings"
)

var mH1 = regexp.MustCompile(`^# (.*)$`)
var mH2 = regexp.MustCompile(`^## (.*)$`)
var mH3 = regexp.MustCompile(`^### (.*)$`)

// replaceHeaders replaces all Gemini headers with HTML headers.
func replaceHeaders(geminiBody string) string {
	lines := strings.Split(geminiBody, "\n")
	for i, line := range lines {
		lines[i] = mH1.ReplaceAllString(line, "<h1>$1</h1>")
	}
	for i, line := range lines {
		lines[i] = mH2.ReplaceAllString(line, "<h2>$1</h2>")
	}
	for i, line := range lines {
		lines[i] = mH3.ReplaceAllString(line, "<h3>$1</h3>")
	}
	return strings.Join(lines, "\n")
}

func ConvertToHTML(geminiBody string) string {
	return replaceHeaders(geminiBody)
}
