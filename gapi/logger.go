package gapi

import (
	"context"
	"net/http"
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
		Str("status_text", statusCode.String()).
		Dur("duration", time.Since(start)).
		Msg("received gRPC")
	return result, err
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler.ServeHTTP(writer, req)
		log.
			Info().
			Str("protocol", "HTTP").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", int(req.Response.StatusCode)).
			Str("status_text", req.Response.Status).
			Dur("duration", time.Since(start)).
			Msg("received gRPC")
	})
}
