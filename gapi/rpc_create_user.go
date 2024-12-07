package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/pauldin91/backend/db/sqlc"
	pb "github.com/pauldin91/backend/pb"
	"github.com/pauldin91/backend/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "method unimplemented")
	}
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user : %s exists", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user : %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}
