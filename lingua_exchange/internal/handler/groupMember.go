package handler

import "github.com/gin-gonic/gin"

var _ GroupMemberHandler = (*groupMemberHandler)(nil)

type GroupMemberHandler interface {
	Members(ctx *gin.Context)
	GetInviteFriends(ctx *gin.Context)
	RemoveMembers(ctx *gin.Context)
	UpdateMemberRemark(ctx *gin.Context)
}

type groupMemberHandler struct{}

func NewGroupMemberHandler() GroupMemberHandler {
	return &groupMemberHandler{}
}

func (g groupMemberHandler) Members(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupMemberHandler) GetInviteFriends(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupMemberHandler) RemoveMembers(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupMemberHandler) UpdateMemberRemark(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}
