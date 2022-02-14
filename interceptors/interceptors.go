package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func ValidateRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Request validation interceptor for method %s. Request: %+v", info.FullMethod, req)
	return handler(ctx, req)
}

func ValidateUser(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// meta, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, errors.New("no metadata in context")
	// }

	// log.Println("Request metadata: ", meta)
	log.Printf("User validation interceptor for method %s. Request: %+v", info.FullMethod, req)
	return handler(ctx, req)
}

func AccessControl(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Access control for method %s. Request: %+v", info.FullMethod, req)
	return handler(ctx, req)
}
