package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "usersvc/api/gen/go/user/v1"
	appuser "usersvc/internal/app/user"
	domain "usersvc/internal/domain/user"
)

type Server struct {
	userv1.UnimplementedUserServiceServer
	CreateUserUC *appuser.CreateUserUseCase
}

func New(createUC *appuser.CreateUserUseCase) *Server {
	return &Server{CreateUserUC: createUC}
}

func (s *Server) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	res, err := s.CreateUserUC.Execute(ctx, appuser.CreateUserCommand{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		switch err {
		case domain.ErrInvalidName, domain.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case appuser.ErrPersist:
			return nil, status.Error(codes.Internal, "failed to persist user")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &userv1.CreateUserResponse{UserId: res.ID}, nil
}
