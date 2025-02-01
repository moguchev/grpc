package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func DebugInterceptor(enabled bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if enabled {
			log.Println("method:", method)
			reqBytes, _ := protojson.Marshal(req.(protoreflect.ProtoMessage))
			log.Println("request:", string(reqBytes))
		}

		rpcErr := invoker(ctx, method, req, reply, cc, opts...)

		if enabled {
			var respBytes []byte
			if m, ok := reply.(protoreflect.ProtoMessage); ok {
				respBytes, _ = protojson.Marshal(m)
				log.Println("repsonse:", string(respBytes))
			}

			if rpcErr != nil {
				log.Println("error:", rpcErr)
			} else {
				log.Println("error:", "<nil>")
			}
		}

		return rpcErr
	}
}
