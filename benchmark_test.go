package gingrpc

import (
	"bytes"
	"github.com/dan-and-dna/gin-grpc/internal/userservice"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RequestForBenchmark(b *testing.B, handler http.Handler, method, path string, body *bytes.Buffer, headers ...Header) {
	// request
	req := httptest.NewRequest(method, path, body)
	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkGinGrpc(b *testing.B) {
	gin.SetMode(gin.TestMode)

	a := &A{}
	option := &AOption{}
	option.SetHandler("/user.userservice/login", &Handler{&userservice.LoginReq{}, a.LoginForBenchmark})
	option.SetHandler("/user.userservice/isauthorized", &Handler{&userservice.IsAuthorizedReq{}, a.IsAuthorized})

	router := gin.New()
	router.POST("/test/:pkg/:service/:method", GinGrpc(option, true))
	body := new(bytes.Buffer)
	body.WriteString(`{"name":"Dan","password":"u12345678"}`)

	RequestForBenchmark(b, router, "POST", "/test/user/userservice/login", body)
}
