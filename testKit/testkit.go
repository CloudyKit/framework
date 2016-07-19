package testKit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewHTTPTester(t testing.TB, a http.Handler) HTTPTester {
	return HTTPTester{T: t, A: a}
}

type HTTPTester struct {
	T testing.TB
	A http.Handler
}

func (appt *HTTPTester) Request(r *http.Request) *Response {
	i := &Response{
		T: appt.T,
		ResponseRecorder: httptest.ResponseRecorder{
			HeaderMap: make(http.Header),
			Body:      new(bytes.Buffer),
			Code:      200,
		},
	}
	appt.A.ServeHTTP(i, r)
	return i
}

func (result *Response) String() string {
	return result.Body.String()
}

func (appt HTTPTester) GetRequest(urlStr string) *Response {
	r, _ := http.NewRequest("GET", urlStr, nil)
	return appt.Request(r)
}

func (r *Response) ExpectOutput(output, errStr string, v ...interface{}) *Response {
	return r.Expect(string(r.Body.Bytes()) == output, errStr, v...)
}

func (r *Response) ExpectOutputContains(output, errStr string, v ...interface{}) *Response {
	return r.Expect(strings.Contains(string(r.Body.Bytes()), output), errStr, v...)
}

func (r *Response) ExpectStatus(status int, errStr string, v ...interface{}) *Response {
	return r.Expect(status == r.Code, errStr, v...)
}

func (r *Response) ExpectRedirect(urlStr, errStr string, v ...interface{}) *Response {
	return r.Expect(r.Header().Get("Location") == urlStr, errStr, v...)
}

func (r *Response) Expect(cond bool, errStr string, v ...interface{}) *Response {
	if !cond {
		r.T.Errorf(errStr, v...)
	}
	return r
}

type Response struct {
	T testing.TB
	httptest.ResponseRecorder
}
