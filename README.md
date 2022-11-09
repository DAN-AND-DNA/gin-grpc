# gin-grpc
The Gin middleware that forwards grpc requests


# [中文文档](./README_cn.md)
1. This middleware enables us to handle requests for different protocols simultaneously with just one piece of code

2. Here is the Grpc middleware [grpc-route](https://github.com/DAN-AND-DNA/grpc-route) 


## Usage Scenarios
- Need to handle both Restful Api and Grpc

- Network Framework and Business Framework (for example: MVC) are separated from each other, just need use the grpc way to write business, such as WebBFF and Services can use one Framework

- Reuse Gin and Grpc community middleware to build your own Microservices framework (Customizable for Observability and Performance)

## Benchmark
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

## Unit Test
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

