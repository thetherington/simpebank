package gapi

import (
	db "github.com/thetherington/simplebank/db/sqlc"
	"github.com/thetherington/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
