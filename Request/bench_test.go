package Request

import (
	"github.com/CloudyKit/framework/App"
	"net/http"
	"testing"
)

type BenchController struct {
	Context
	*testing.B
}

func invokeNextMiddleware(rf Channel) {
	rf.Next()
}

func (bb *BenchController) Mux(m App.Mapper) {
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
	context := c.Child()
	defer context.Done()
}

func (c *BenchController) SimpleHandler() {
	if c.B == nil {
		panic("Can't load b from Di.Context")
	}
}

var benchApp = App.New()

func init() {
	benchApp.AddController(
		&BenchController{},
	)
}
func BenchmarkFlowRequest(b *testing.B) {
	benchApp.Put(b)
	request, _ := http.NewRequest("GET", "/", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}

func BenchmarkFlowRequestMiddleware1(b *testing.B) {
	request, _ := http.NewRequest("GET", "/middlewares/1", nil)
	benchApp.Put(b)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}

func BenchmarkFlowRequestMiddleware4(b *testing.B) {
	benchApp.Put(b)
	request, _ := http.NewRequest("GET", "/middlewares/4", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}
func BenchmarkFlowRequestMiddleware24(b *testing.B) {
	benchApp.Put(b)
	request, _ := http.NewRequest("GET", "/middlewares/24", nil)
	for i := 0; i < b.N; i++ {
		benchApp.Router.ServeHTTP(nil, request)
	}
}
