package sso

import (
	ssov1 "ecodepost/pb/sso/v1"
)

type GrpcServer struct{}

var _ ssov1.SsoServer = &GrpcServer{}
