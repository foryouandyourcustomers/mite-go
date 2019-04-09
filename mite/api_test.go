package mite_test

const testApiKey = "f00bar"
const testClientVersion = "vX"
const testUserAgent = "mite-go/" + testClientVersion + " (+github.com/leanovate/mite-go)"

type recorder struct {
	method      string
	url         string
	contentType string
	userAgent   string
	miteKey     string
}
