package gingrpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type PathToServiceName func(c *gin.Context) string
type GrpcHandlers map[string]Handler

type Handler struct {
	ReqProto proto.Message
	Handle   func(context.Context, interface{}) (interface{}, error)
}

type GrpcCtxOption interface {
	Apply(ctx context.Context) context.Context
}

func GinGrpc(path2Service PathToServiceName, grpcHandlers GrpcHandlers, httpHeader bool, options ...GrpcCtxOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 拿协议
		key := path2Service(c)

		// 填充协议
		bodyBuffer, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": codes.Internal, "error_desc": codes.Internal.String(), "message": err.Error()})
			return
		}

		handler, ok := grpcHandlers[key]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": codes.InvalidArgument, "error_desc": codes.InvalidArgument.String(), "message": "no such proto"})
			return
		}

		reqProto := handler.ReqProto
		err = protojson.Unmarshal(bodyBuffer, reqProto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": codes.InvalidArgument, "error_desc": codes.InvalidArgument.String(), "message": "bad json"})
			return
		}

		c.Set("gin_grpc_req", reqProto)
		c.Set("gin_grpc_handler", handler)

		// 处理协议
		handleGrpcRequest(httpHeader, options...)(c)

		// 返回结果
		rawErr, ok := c.Get("gin_grpc_err")
		if ok {
			s := status.Convert(rawErr.(error))
			c.JSON(http.StatusOK, gin.H{"code": s.Code(), "error_desc": s.Code().String(), "message": s.Message()})
			return
		}
		rawResp, _ := c.Get("gin_grpc_resp")
		c.JSON(http.StatusOK, rawResp)
		return
	}
}

func handleGrpcRequest(httpHeader bool, options ...GrpcCtxOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqProto, ok := c.Get("gin_grpc_req")
		if !ok {
			return
		}

		rawHandler, ok := c.Get("gin_grpc_handler")
		if !ok {
			return
		}

		var ctx context.Context = c

		// 拿http头部
		if httpHeader {
			var md metadata.MD = metadata.MD{}
			for key, val := range c.Request.Header {
				md.Append(key, val...)
			}
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		// 比如zap日志
		for _, option := range options {
			ctx = option.Apply(ctx)
		}

		respProto, err := rawHandler.(Handler).Handle(ctx, reqProto)
		if err != nil {
			c.Set("gin_grpc_err", err)
			return
		}
		c.Set("gin_grpc_resp", respProto)
	}
}
