package webproxys

import (
	"fmt"
	"regexp"
	"strings"
)

var mBlank = regexp.MustCompile(`^\s*$`)
var mH1 = regexp.MustCompile(`^# (.*)$`)
var mH2 = regexp.MustCompile(`^## (.*)$`)
var mH3 = regexp.MustCompile(`^### (.*)$`)

// =>[<whitespace>]<URL>[<whitespace><USER-FRIENDLY LINK NAME>]
var mLink = regexp.MustCompile(`^=>\s*(\S+)\s*(.*)$`)

type replacer = func(string) *string

// replaceByLine performs text replacement on a line-by-line basis.
func replaceByLine(geminiBody string, replacers ...replacer) string {
	lines := strings.Split(geminiBody, "\n")
	for _, replacer := range replacers {
		for i, line := range lines {
			out := replacer(line)
			if out != nil {
				lines[i] = *out
			}
		}
	}
	return strings.Join(lines, "\n")
}

func replacerForRegexp(matcher *regexp.Regexp, replacement string) replacer {
	return func(line string) *string {
		if matcher.MatchString(line) {
			val := matcher.ReplaceAllString(line, replacement)
			return &val
		}
		return nil
	}
}

func linkReplacer(line string) *string {
	match := mLink.FindStringSubmatch(line)
	if match != nil {
		url, desc := match[1], match[2]
		if desc == "" {
			desc = url
		}
		a := fmt.Sprintf("<a href=\"%s\">%s</a><br>", url, desc)
		return &a
	}
	return nil
}

func ConvertToHTML(geminiBody string) string {
	body := strings.TrimSpace(geminiBody)
	return replaceByLine(
		body,
		replacerForRegexp(mH1, "<h1>$1</h1>"),
		replacerForRegexp(mH2, "<h2>$1</h2>"),
		replacerForRegexp(mH3, "<h3>$1</h3>"),
		linkReplacer,
		replacerForRegexp(mBlank, "<br>"),
	)
}
