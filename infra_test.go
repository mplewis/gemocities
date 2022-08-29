package gemocities_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	neturl "net/url"
	"regexp"

	"git.sr.ht/~adnano/go-gemini"
)

var clientCerts *tls.Certificate
var linkMatcher = regexp.MustCompile(`(?m)^=>\s*([^\s]+)\s+(.+)$`)

type ResponseBuffer struct {
	Data      []byte
	MediaType string
	Status    gemini.Status
	Meta      string
}

func (r *ResponseBuffer) Write(b []byte) (int, error) {
	r.Data = append(r.Data, b...)
	return len(b), nil
}

func (r *ResponseBuffer) WriteHeader(status gemini.Status, meta string) {
	r.Status = status
	r.Meta = meta
}

func (r *ResponseBuffer) SetMediaType(mediatype string) {
	r.MediaType = mediatype
}

func (r *ResponseBuffer) Body() string {
	return string(r.Data)
}

func (r *ResponseBuffer) Flush() error { return nil }

func (r *ResponseBuffer) Links() Links {
	var links Links
	matches := linkMatcher.FindAllStringSubmatch(r.Body(), -1)
	for _, match := range matches {
		links = append(links, Link{URL: match[1], Text: match[2]})
	}
	return links
}

type Links []Link

func (l Links) WithText(text string) (Link, bool) {
	for _, link := range l {
		if link.Text == text {
			return link, true
		}
	}
	return Link{}, false
}

type Link struct {
	URL  string
	Text string
}

type Requestor struct {
	*gemini.Server
}

func (r Requestor) Request(url string, cert *tls.Certificate) ResponseBuffer {
	u, err := neturl.Parse(url)
	if err != nil {
		panic(err)
	}
	req := gemini.Request{URL: u, Certificate: cert}
	var resp ResponseBuffer
	r.Server.Handler.ServeGemini(context.Background(), &resp, &req)
	return resp
}

func (r Requestor) RequestInput(url string, cert *tls.Certificate, input string) ResponseBuffer {
	return r.Request(fmt.Sprintf("%s?%s", url, input), cert)
}

func ClientCerts() *tls.Certificate {
	if clientCerts == nil {
		cert, err := tls.LoadX509KeyPair("test/certs/test_user.crt", "test/certs/test_user.key")
		if err != nil {
			panic(err)
		}
		raw := cert.Certificate[0]
		cert.Leaf, err = x509.ParseCertificate(raw)
		if err != nil {
			panic(err)
		}
		clientCerts = &cert
	}
	return clientCerts
}
