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

package common

import (
	"fmt"
	"testing"
)

func TestGenQS(t *testing.T) {

	generatedURL := GenQS(nil, "/users")("username", "myusername")
	expectedURL := "/users?username=myusername"
	if generatedURL != expectedURL {
		t.Errorf("want %q got %q", expectedURL, generatedURL)
	}
}

func BenchmarkGenQS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generatedURL := GenQS(nil, "/users")("username", "myusername")
		expectedURL := "/users?username=myusername"
		if generatedURL != expectedURL {
			b.Errorf("want %q got %q", expectedURL, generatedURL)
		}
	}
}

func BenchmarkGenURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generatedURL := GenURL(nil, "/users") + fmt.Sprintf("?%s=%s", "username", "myusername")
		expectedURL := "/users?username=myusername"
		if generatedURL != expectedURL {
			b.Errorf("want %q got %q", expectedURL, generatedURL)
		}
	}
}
