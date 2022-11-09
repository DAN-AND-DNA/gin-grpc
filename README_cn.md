# gin-grpc
1. 这个中间件使得我们可以只用一份代码同时处理不同协议的请求

2. 利用gin的中间件，以grpc的方式处理restful api请求

## 使用场景
- 需要同时处理restful api 和 grpc

- 网络框架和业务框架（比如MVC）相互分离，只需要使用grpc的方式写业务，比如WebBFF和后排Services统一框架

- 复用gin和grpc社区中间件，构建自己的微服务框架（定制化的可观测性和性能）

## 性能测试
```
goos: windows
goarch: amd64
pkg: github.com/dan-and-dna/gin-grpc
cpu: Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz
BenchmarkGinGrpc
BenchmarkGinGrpc-12      3528080              1675 ns/op            1496 B/op
              13 allocs/op
PASS
```

## 单元测试

```
=== RUN   TestGinGrpc
=== RUN   TestGinGrpc/TestJsonUnmarshalAndMarshal
=== RUN   TestGinGrpc/TestReturnNil
=== RUN   TestGinGrpc/TestRequestError
=== RUN   TestGinGrpc/TestRightRequest
=== RUN   TestGinGrpc/TestOmitempty
=== RUN   TestGinGrpc/TestBadBody
=== RUN   TestGinGrpc/TestEmptyBody
=== RUN   TestGinGrpc/TestBadPath
--- PASS: TestGinGrpc (0.00s)
    --- PASS: TestGinGrpc/TestJsonUnmarshalAndMarshal (0.00s)
    --- PASS: TestGinGrpc/TestReturnNil (0.00s)
    --- PASS: TestGinGrpc/TestRequestError (0.00s)
    --- PASS: TestGinGrpc/TestRightRequest (0.00s)
    --- PASS: TestGinGrpc/TestOmitempty (0.00s)
    --- PASS: TestGinGrpc/TestBadBody (0.00s)
    --- PASS: TestGinGrpc/TestEmptyBody (0.00s)
    --- PASS: TestGinGrpc/TestBadPath (0.00s)
PASS
```