package util

import (
	"google.golang.org/grpc"
)

func NewClientConnection(target string) func(...grpc.DialOption) (*grpc.ClientConn, error) {
	return func(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return grpc.Dial(target, opts...)
	}
}
