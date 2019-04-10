package mite_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	testApiKey        = "key"
	testClientVersion = "test"
	testUserAgent     = "mite-go/" + testClientVersion + " (+github.com/leanovate/mite-go)"
)

type Recorder struct {
	reqMethod string
	reqUri    string
	reqBody   []byte
	reqHeader http.Header
	resHeader http.Header
	resBody   string
	resStatus int
}

func NewRecorder() *Recorder {
	return &Recorder{reqHeader: make(http.Header), resHeader: make(http.Header), resStatus: 200}
}

func (r *Recorder) RequestMethod() string {
	return r.reqMethod
}

func (r *Recorder) RequestURI() string {
	return r.reqUri
}

func (r *Recorder) RequestContentType() string {
	return r.RequestHeader("Content-Type")
}

func (r *Recorder) RequestUserAgent() string {
	return r.RequestHeader("User-Agent")
}

func (r *Recorder) RequestMiteKey() string {
	return r.RequestHeader("X-MiteApiKey")
}

func (r *Recorder) RequestHeader(header string) string {
	return r.reqHeader.Get(header)
}

func (r *Recorder) RequestBody() []byte {
	return r.reqBody
}

func (r *Recorder) RequestBodyCanonical() string {
	return string(canonicalJSON(r.reqBody))
}

func (r *Recorder) ResponseContentType(contentType string) *Recorder {
	r.resHeader.Add("Content-Type", contentType)
	return r
}

func (r *Recorder) ResponseLocation(location string) *Recorder {
	r.resHeader.Add("Location", location)
	return r
}

func (r *Recorder) ResponseBody(body string) *Recorder {
	r.resBody = body
	return r
}

func (r *Recorder) ResponseStatus(status int) *Recorder {
	r.resStatus = status
	return r
}

func (r *Recorder) Handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		r.reqMethod = request.Method
		r.reqUri = request.RequestURI
		r.reqHeader = request.Header
		r.reqBody = b

		for k, vs := range r.resHeader {
			for _, v := range vs {
				writer.Header().Add(k, v)
			}
		}

		writer.WriteHeader(r.resStatus)

		if r.resBody == "" {
			return
		}

		_, err = writer.Write([]byte(r.resBody))
		if err != nil {
			panic(err)
		}
	})
}

func CanonicalString(s string) string {
	return string(canonicalJSON([]byte(s)))
}

func canonicalJSON(b []byte) []byte {
	buf := &bytes.Buffer{}
	err := json.Indent(buf, b, "", "  ")
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
