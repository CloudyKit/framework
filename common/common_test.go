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
