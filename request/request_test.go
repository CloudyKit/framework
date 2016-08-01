package request

import (
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/router"
	"testing"
)

func TestContext_Advance(t *testing.T) {
	c := new(Context)

	counter := 0
	handler := HandlerFunc(func(c *Context) {
		counter++
		c.Advance()
	})

	Advance(c, "TestHandler", nil, nil, router.Parameter{}, scope.New(), []Handler{
		handler,
		handler,
		handler,
		handler,
		handler,
	})

	if counter != 5 {
		t.Errorf("Not all handlers executed: want 5 got %v", counter)
	}
}
