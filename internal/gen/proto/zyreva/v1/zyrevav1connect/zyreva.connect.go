// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/zyreva/v1/zyreva.proto

package zyrevav1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/zyreva/zyreva/internal/gen/proto/zyreva/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// AppServiceName is the fully-qualified name of the AppService service.
	AppServiceName = "zyreva.v1.AppService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AppServiceAnalyzeProcedure is the fully-qualified name of the AppService's Analyze RPC.
	AppServiceAnalyzeProcedure = "/zyreva.v1.AppService/Analyze"
)

// AppServiceClient is a client for the zyreva.v1.AppService service.
type AppServiceClient interface {
	Analyze(context.Context, *connect_go.Request[v1.AnalyzeRequest]) (*connect_go.ServerStreamForClient[v1.AnalyzeResponse], error)
}

// NewAppServiceClient constructs a client for the zyreva.v1.AppService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAppServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AppServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &appServiceClient{
		analyze: connect_go.NewClient[v1.AnalyzeRequest, v1.AnalyzeResponse](
			httpClient,
			baseURL+AppServiceAnalyzeProcedure,
			opts...,
		),
	}
}

// appServiceClient implements AppServiceClient.
type appServiceClient struct {
	analyze *connect_go.Client[v1.AnalyzeRequest, v1.AnalyzeResponse]
}

// Analyze calls zyreva.v1.AppService.Analyze.
func (c *appServiceClient) Analyze(ctx context.Context, req *connect_go.Request[v1.AnalyzeRequest]) (*connect_go.ServerStreamForClient[v1.AnalyzeResponse], error) {
	return c.analyze.CallServerStream(ctx, req)
}

// AppServiceHandler is an implementation of the zyreva.v1.AppService service.
type AppServiceHandler interface {
	Analyze(context.Context, *connect_go.Request[v1.AnalyzeRequest], *connect_go.ServerStream[v1.AnalyzeResponse]) error
}

// NewAppServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAppServiceHandler(svc AppServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	appServiceAnalyzeHandler := connect_go.NewServerStreamHandler(
		AppServiceAnalyzeProcedure,
		svc.Analyze,
		opts...,
	)
	return "/zyreva.v1.AppService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AppServiceAnalyzeProcedure:
			appServiceAnalyzeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAppServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAppServiceHandler struct{}

func (UnimplementedAppServiceHandler) Analyze(context.Context, *connect_go.Request[v1.AnalyzeRequest], *connect_go.ServerStream[v1.AnalyzeResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("zyreva.v1.AppService.Analyze is not implemented"))
}
