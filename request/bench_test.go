package request_test

import (
	"github.com/CloudyKit/framework/app"
	. "github.com/CloudyKit/framework/request"
	"net/http"
	"testing"
)

type BenchController struct {
	*Context
	*testing.B
}

func invokeNextMiddleware(rf ContextChain) {
	rf.Next()
}

func (bb *BenchController) Mux(m app.Mapper) {
	m.AddHandler("GET", "/", "SimpleHandler")
	m.AddHandler("GET", "/middlewares/1", "SimpleHandler", invokeNextMiddleware)
	m.AddHandler("GET", "/middlewares/4", "SimpleHandler", invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware)
	m.AddHandler("GET", "/middlewares/24", "SimpleHandler",
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
		invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware, invokeNextMiddleware,
	)
}

func (c *BenchController) RegisterHandler() {
	context := c.Di.Child()
	defer context.Done()
}

func (c *BenchController) SimpleHandler() {
	if c.B == nil {
		panic("Can't load b from Di.Context")
	}
}

var benchApp = app.New()

func init() {
	benchApp.AddController(
		&BenchController{},
	)
}
func BenchmarkFlowRequest(b *testing.B) {
	benchApp.Di.Map(b)
	request, _ := http.NewRequest("GET", "/", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}

func BenchmarkFlowRequestMiddleware1(b *testing.B) {
	request, _ := http.NewRequest("GET", "/middlewares/1", nil)
	benchApp.Di.Map(b)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}

func BenchmarkFlowRequestMiddleware4(b *testing.B) {
	benchApp.Di.Map(b)
	request, _ := http.NewRequest("GET", "/middlewares/4", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}
func BenchmarkFlowRequestMiddleware24(b *testing.B) {
	benchApp.Di.Map(b)
	request, _ := http.NewRequest("GET", "/middlewares/24", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}
