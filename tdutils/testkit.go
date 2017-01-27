// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tdutils

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
