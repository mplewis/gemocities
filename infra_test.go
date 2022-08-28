package gemocities_test

import (
	"context"
	"crypto/tls"
	neturl "net/url"

	"git.sr.ht/~adnano/go-gemini"
)

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

func Request(srv *gemini.Server, url string, cert *tls.Certificate) ResponseBuffer {
	u, err := neturl.Parse(url)
	if err != nil {
		panic(err)
	}
	req := gemini.Request{URL: u, Certificate: cert}
	var resp ResponseBuffer
	srv.Handler.ServeGemini(context.Background(), &resp, &req)
	return resp
}
