package webproxys

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var mH = regexp.MustCompile(`^(#{1,3}) (.*)$`)
var mLink = regexp.MustCompile(`^=>\s*(\S+)\s*(.*)$`)

const markPre = "```"

type TaggedChunk struct {
	Tag   string
	Body  string
	Match []string
}

func process(chunks []TaggedChunk) string {
	out := ""
	for _, chunk := range chunks {
		switch chunk.Tag {
		case "header":
			hLevelRaw, text := chunk.Match[1], chunk.Match[2]
			hLevel := "h1"
			if hLevelRaw == "##" {
				hLevel = "h2"
			} else if hLevelRaw == "###" {
				hLevel = "h3"
			}
			out += fmt.Sprintf("<%s>%s</%s>\n", hLevel, html.EscapeString(text), hLevel)
		case "link":
			href, text := chunk.Match[1], chunk.Match[2]
			if text == "" {
				text = href
			}
			// TODO: Escape href
			out += fmt.Sprintf("<p><a href=\"%s\">%s</a></p>\n", href, html.EscapeString(text))
		case "preformatted":
			out += fmt.Sprintf("<p><pre>%s</pre></p>\n", html.EscapeString(chunk.Body))
		default:
			out += fmt.Sprintf("<p>%s</p>\n", html.EscapeString(chunk.Body))
		}
	}
	return out
}

func ConvertToHTML(geminiBody string) string {
	chunks := []TaggedChunk{}
	chunk := []string{}
	pre := false

	for _, c := range strings.Split(strings.TrimSpace(geminiBody), "\n") {
		if c == markPre {
			pre = !pre
			if !pre {
				chunks = append(chunks, TaggedChunk{"preformatted", strings.Join(chunk, "\n"), nil})
				chunk = []string{}
			}
			continue
		}
		if pre {
			chunk = append(chunk, c)
			continue
		}

		if m := mH.FindStringSubmatch(c); m != nil {
			chunks = append(chunks, TaggedChunk{"header", c, m})
			continue
		}
		if m := mLink.FindStringSubmatch(c); m != nil {
			chunks = append(chunks, TaggedChunk{"link", c, m})
			continue
		}
		chunks = append(chunks, TaggedChunk{"", c, nil})
	}

	return process(chunks)
}
