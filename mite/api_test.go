package mite_test

import (
	"bytes"
	"encoding/json"
)

const (
	testApiKey        = "key"
	testClientVersion = "test"
	testUserAgent     = "mite-go/" + testClientVersion + " (+github.com/leanovate/mite-go)"
)

type recorder struct {
	method      string
	url         string
	body        []byte
	contentType string
	userAgent   string
	miteKey     string
}

func prettifyJson(b []byte, indent string) []byte {
	buf := &bytes.Buffer{}
	err := json.Indent(buf, b, "", indent)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
