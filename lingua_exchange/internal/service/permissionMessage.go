package service

import (
	"context"

	"lingua_exchange/internal/types"
)

// 消息权限控制类
var _ IPermissionService = (*PermissionService)(nil)

type IPermissionService interface {
	IsAuth(ctx context.Context, opt *types.AuthOption) error
}

type PermissionService struct {
}

func (p PermissionService) IsAuth(ctx context.Context, opt *types.AuthOption) error {
	// TODO implement me
	panic("implement me")
}
