package main

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func WithProtoValidate(validator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if protoMsg, ok := req.(proto.Message); ok {
			if err := validator.Validate(protoMsg); err != nil {
				st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())
				st, _ = st.WithDetails(&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "request",
							Description: err.Error(),
						},
					},
				})

				return nil, st.Err()
			}
		}

		return handler(ctx, err)
	}
}
