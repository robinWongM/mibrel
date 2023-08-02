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
	// GitServiceName is the fully-qualified name of the GitService service.
	GitServiceName = "zyreva.v1.GitService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GitServiceCloneProcedure is the fully-qualified name of the GitService's Clone RPC.
	GitServiceCloneProcedure = "/zyreva.v1.GitService/Clone"
)

// GitServiceClient is a client for the zyreva.v1.GitService service.
type GitServiceClient interface {
	Clone(context.Context, *connect_go.Request[v1.CloneRequest]) (*connect_go.Response[v1.CloneResponse], error)
}

// NewGitServiceClient constructs a client for the zyreva.v1.GitService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGitServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) GitServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &gitServiceClient{
		clone: connect_go.NewClient[v1.CloneRequest, v1.CloneResponse](
			httpClient,
			baseURL+GitServiceCloneProcedure,
			opts...,
		),
	}
}

// gitServiceClient implements GitServiceClient.
type gitServiceClient struct {
	clone *connect_go.Client[v1.CloneRequest, v1.CloneResponse]
}

// Clone calls zyreva.v1.GitService.Clone.
func (c *gitServiceClient) Clone(ctx context.Context, req *connect_go.Request[v1.CloneRequest]) (*connect_go.Response[v1.CloneResponse], error) {
	return c.clone.CallUnary(ctx, req)
}

// GitServiceHandler is an implementation of the zyreva.v1.GitService service.
type GitServiceHandler interface {
	Clone(context.Context, *connect_go.Request[v1.CloneRequest]) (*connect_go.Response[v1.CloneResponse], error)
}

// NewGitServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGitServiceHandler(svc GitServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	gitServiceCloneHandler := connect_go.NewUnaryHandler(
		GitServiceCloneProcedure,
		svc.Clone,
		opts...,
	)
	return "/zyreva.v1.GitService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GitServiceCloneProcedure:
			gitServiceCloneHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGitServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGitServiceHandler struct{}

func (UnimplementedGitServiceHandler) Clone(context.Context, *connect_go.Request[v1.CloneRequest]) (*connect_go.Response[v1.CloneResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("zyreva.v1.GitService.Clone is not implemented"))
}
