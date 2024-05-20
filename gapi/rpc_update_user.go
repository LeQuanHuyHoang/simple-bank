package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/utils"
	"simple-bank/val"
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invaildArgumentError(violations)
	}

	args := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.GetFullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.GetEmail != nil,
		},
	}

	if req.Password != nil {
		hasdedPassword, err := utils.HasdPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		args.HashedPassword = sql.NullString{
			String: hasdedPassword,
			Valid:  true,
		}
		args.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsq := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	fmt.Println(user)

	return rsq, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		validations = append(validations, fieldViolation("username", err))
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			validations = append(validations, fieldViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			validations = append(validations, fieldViolation("email", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			validations = append(validations, fieldViolation("full_name", err))
		}
	}

	return validations
}
