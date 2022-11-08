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

type Handler struct {
	GetProto    func() proto.Message
	PutProto    func(proto.Message)
	HandleProto func(context.Context, interface{}) (interface{}, error)
}

type Option interface {
	PathToGrpcService(c *gin.Context) string
	GetHandler(string) (*Handler, bool)
	SetHandler(string, *Handler)
}

type GrpcCtxOption interface {
	Apply(ctx context.Context) context.Context
}

func GinGrpc(option Option, httpHeader bool, grpcCtxOptions ...GrpcCtxOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if option == nil {
			return
		}
		// 拿协议
		key := option.PathToGrpcService(c)

		// 填充协议
		bodyBuffer, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": codes.Internal, "error_desc": codes.Internal.String(), "message": err.Error()})
			return
		}

		handler, ok := option.GetHandler(key)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": codes.InvalidArgument, "error_desc": codes.InvalidArgument.String(), "message": "unknown request"})
			return
		}

		if handler == nil || handler.GetProto == nil || handler.HandleProto == nil {
			return
		}

		reqProto := handler.GetProto()
		err = protojson.Unmarshal(bodyBuffer, reqProto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": codes.InvalidArgument, "error_desc": codes.InvalidArgument.String(), "message": "bad json"})
			return
		}

		c.Set("gin_grpc_req", reqProto)
		c.Set("gin_grpc_handler", handler)

		// 处理协议
		handleGrpcRequest(httpHeader, grpcCtxOptions...)(c)

		if handler.PutProto != nil {
			handler.PutProto(reqProto)
		}

		// 返回结果
		rawErr, ok := c.Get("gin_grpc_err")
		if ok {
			s := status.Convert(rawErr.(error))
			if s.Code() == codes.Internal {
				c.JSON(http.StatusInternalServerError, gin.H{"code": s.Code(), "error_desc": s.Code().String(), "message": s.Message()})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": s.Code(), "error_desc": s.Code().String(), "message": s.Message()})
			}
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

		respProto, err := rawHandler.(*Handler).HandleProto(ctx, reqProto)
		if err != nil {
			c.Set("gin_grpc_err", err)
			return
		}
		c.Set("gin_grpc_resp", respProto)
	}
}
