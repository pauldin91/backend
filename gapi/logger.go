package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	result, err := handler(ctx, req)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	log.
		Info().
		Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_test", statusCode.String()).
		Dur("duration", time.Since(start)).
		Msg("received gRPC")
	return result, err
}
