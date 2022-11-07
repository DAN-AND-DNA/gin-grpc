# gin-grpc

利用gin的中间件，以grpc的方式处理restful api请求

## 使用场景
- 需要同时处理restful 和 grpc
- 网络框架和业务框架相互分离，只使用grpc的方式写业务，比如WebBFF和后排Services统一框架
- 复用gin和grpc社区中间件，构建自己的微服务框架（定制化的可观测性和性能）
