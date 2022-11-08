package gingrpc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dan-and-dna/gin-grpc/internal/userservice"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Header struct {
	Key   string
	Value string
}

func Request(handler http.Handler, method, path string, body *bytes.Buffer, headers ...Header) *httptest.ResponseRecorder {
	// request
	req := httptest.NewRequest(method, path, body)
	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	// result
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

type A struct {
}

func (a *A) LoginReq() proto.Message {
	return &userservice.LoginReq{}
}

func (a *A) IsAuthorizedReq() proto.Message {
	return &userservice.IsAuthorizedReq{}
}

func (a *A) Login(ctx context.Context, req interface{}) (interface{}, error) {
	reqProto, ok := req.(*userservice.LoginReq)
	if !ok {
		return nil, status.Error(codes.Internal, "bad proto")
	}

	if reqProto.GetName() == "Dan" && reqProto.GetPassword() == "u12345678" {
		resp := userservice.LoginResp{
			Token: "token1234567",
			UserInfo: &userservice.UserInfo{
				Uid:      "10001",
				Username: "Dan",
				Age:      99,
			},
		}
		return resp, nil
	}

	return nil, status.Error(codes.InvalidArgument, "invalid user or password")
}

func (a *A) IsAuthorized(ctx context.Context, req interface{}) (interface{}, error) {
	reqProto, ok := req.(*userservice.IsAuthorizedReq)
	if !ok {
		return nil, status.Error(codes.Internal, "bad proto")
	}

	if reqProto.GetToken() == "token1234567" {
		return &userservice.IsAuthorizedResp{
			IsAuthorized: true,
		}, nil
	}

	return &userservice.IsAuthorizedResp{
		IsAuthorized: false,
	}, nil
}

type AOption struct {
	handlers map[string]*Handler
}

func (a *AOption) PathToGrpcService(c *gin.Context) string {
	pkg := c.Param("pkg")
	service := c.Param("service")
	method := c.Param("method")

	return fmt.Sprintf("/%s.%s/%s", pkg, service, method)
}

func (a *AOption) GetHandler(key string) (*Handler, bool) {
	if h, ok := a.handlers[key]; ok {
		return h, true
	}

	return nil, false
}

func (option *AOption) SetHandler(key string, handler *Handler) {
	if option.handlers == nil {
		option.handlers = map[string]*Handler{}
	}

	option.handlers[key] = handler
}

func TestGinGrpc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	a := &A{}
	option := &AOption{}
	option.SetHandler("/user.userservice/login", &Handler{a.LoginReq, nil, a.Login})
	option.SetHandler("/user.userservice/isauthorized", &Handler{a.IsAuthorizedReq, nil, a.IsAuthorized})

	router := gin.New()
	router.POST("/test/:pkg/:service/:method", GinGrpc(option, true))

	type args struct {
		body string
		path string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantCode int
	}{
		{
			name:     "TestJsonUnmarshalAndMarshal",
			args:     args{body: `{"name":"Dan","password":"u12345678"}`, path: "/test/user/userservice/login"},
			want:     `{"token":"token1234567","user_info":{"uid":"10001","username":"Dan","age":99}}`,
			wantCode: http.StatusOK,
		},
		{
			name:     "TestRequestError",
			args:     args{body: `{"name":"Dan","password":"12345678"}`, path: "/test/user/userservice/login"},
			want:     `{"code":3,"error_desc":"InvalidArgument","message":"invalid user or password"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "TestRightRequest",
			args:     args{body: `{"token":"token1234567"}`, path: "/test/user/userservice/isauthorized"},
			want:     `{"is_authorized":true}`,
			wantCode: http.StatusOK,
		},
		{
			name:     "TestOmitempty",
			args:     args{body: `{"token":"token12345678"}`, path: "/test/user/userservice/isauthorized"},
			want:     `{}`,
			wantCode: http.StatusOK,
		},
		{
			name:     "TestBadBody",
			args:     args{body: `{"name":"Dan","age":30}`, path: "/test/user/userservice/login"},
			want:     `{"code":3,"error_desc":"InvalidArgument","message":"bad json"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "TestEmptyBody",
			args:     args{body: `{}`, path: "/test/user/userservice/login"},
			want:     `{"code":3,"error_desc":"InvalidArgument","message":"invalid user or password"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "TestBadPath",
			args:     args{body: `{"name":"Dan","age":30}`, path: "/test/user/userservice/getusers"},
			want:     `{"code":3,"error_desc":"InvalidArgument","message":"unknown request"}`,
			wantCode: http.StatusBadRequest,
		},
	}
	body := new(bytes.Buffer)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body.Reset()
			body.WriteString(tt.args.body)

			w := Request(router, "POST", tt.args.path, body, Header{Key: "Content-Type", Value: "application/json"})
			resp := w.Result()
			bodyBuffer, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.want, string(bodyBuffer))
		})
	}
}
