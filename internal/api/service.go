package api

import (
	desc "user/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserServiceV1Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}
