package webproxys

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var mBlank = regexp.MustCompile(`^\s*$`)
var mH1 = regexp.MustCompile(`^# (.*)$`)
var mH2 = regexp.MustCompile(`^## (.*)$`)
var mH3 = regexp.MustCompile(`^### (.*)$`)
var mLink = regexp.MustCompile(`^=>\s*(\S+)\s*(.*)$`)

const markPre = "```"

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

// replacerForRegexp replaces matching regexps with the given replacement string.
func replacerForRegexp(matcher *regexp.Regexp, replacement string) replacer {
	return func(line string) *string {
		if matcher.MatchString(line) {
			val := matcher.ReplaceAllString(line, replacement)
			return &val
		}
		return nil
	}
}

// linkReplacer replaces Gemini links with HTML links.
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

func preReplace(geminiBody string) string {
	lines := strings.Split(geminiBody, "\n")
	var pre bool
	for i, line := range lines {
		if line == markPre {
			pre = !pre
			if pre {
				lines[i] = "<pre>"
			} else {
				lines[i] = "</pre>"
			}
		} else if pre {
			lines[i] = html.EscapeString(line)
		}
	}
	return strings.Join(lines, "\n")
}

func chunkByPre(geminiBody string) ([]string, []string) {
	var normal []string
	var preformatted []string
	var pre bool
	var chunk []string
	done := func(body string) {
		if pre {
			normal = append(normal, body)
		} else {
			preformatted = append(preformatted, body)
		}
	}

	for _, line := range strings.Split(geminiBody, "\n") {
		if line == markPre {
			pre = !pre
			done(strings.Join(chunk, "\n"))
			chunk = []string{}
		} else {
			chunk = append(chunk, line)
		}
	}
	if len(chunk) > 0 {
		done(strings.Join(chunk, "\n"))
	}
	return normal, preformatted
}

func ConvertToHTML(geminiBody string) string {
	// TODO: Escape HTML in the body.
	// Maybe handle preformatted blocks separately?
	body := strings.TrimSpace(geminiBody)
	normal, pre := chunkByPre(body)
	processed := []string{}
	for i, nc := range normal {
		processed = append(processed, replaceByLine(
			nc,
			linkReplacer,
			replacerForRegexp(mH1, "<h1>$1</h1>"),
			replacerForRegexp(mH2, "<h2>$1</h2>"),
			replacerForRegexp(mH3, "<h3>$1</h3>"),
			replacerForRegexp(mBlank, "<br>"),
		))
		if i < len(pre) {
			processed = append(processed, html.EscapeString(pre[i]))
		}
	}
	body = strings.Join(processed, "\n")
	return body
}
