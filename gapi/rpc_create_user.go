package gapi

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/utils"
)

func (server Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hasdedPassword, err := utils.HasdPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	args := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hasdedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username already exist: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsq := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsq, nil
}
