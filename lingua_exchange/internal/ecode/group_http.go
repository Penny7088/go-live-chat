package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// group business-level http error codes.
// the groupNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	groupNO       = 60
	groupName     = "group"
	groupBaseCode = errcode.HCode(groupNO)

	ErrCreateGroup     = errcode.NewError(groupBaseCode+1, "failed to create "+groupName)
	ErrDeleteByIDGroup = errcode.NewError(groupBaseCode+2, "failed to delete "+groupName)
	ErrUpdateByIDGroup = errcode.NewError(groupBaseCode+3, "failed to update "+groupName)
	ErrGetByIDGroup    = errcode.NewError(groupBaseCode+4, "failed to get "+groupName+" details")
	ErrListGroup       = errcode.NewError(groupBaseCode+5, "failed to list of "+groupName)

	ErrDeleteByIDsGroup             = errcode.NewError(groupBaseCode+6, "failed to delete by batch ids "+groupName)
	ErrGetByConditionGroup          = errcode.NewError(groupBaseCode+7, "failed to get "+groupName+" details by conditions")
	ErrListByIDsGroup               = errcode.NewError(groupBaseCode+8, "failed to list by batch ids "+groupName)
	ErrListByLastIDGroup            = errcode.NewError(groupBaseCode+9, "failed to list by last id "+groupName)
	ErrGroupDismiss                 = errcode.NewError(groupBaseCode+10, "group is dismiss "+groupName)
	ErrGroupAlreadyDismiss          = errcode.NewError(groupBaseCode+11, "group already is  dismiss "+groupName)
	ErrGroupInviteFriendsNil        = errcode.NewError(groupBaseCode+12, "group invite friends is null "+groupName)
	ErrGroupInviteNotPermission     = errcode.NewError(groupBaseCode+13, "group invite friends is not permission "+groupName)
	ErrGroupInviteFailed            = errcode.NewError(groupBaseCode+14, "group invite friends failed "+groupName)
	ErrGroupNotExist                = errcode.NewError(groupBaseCode+15, "group not exist  "+groupName)
	ErrGroupNotPermission           = errcode.NewError(groupBaseCode+16, "group not permission  "+groupName)
	ErrGroupDetailsFailed           = errcode.NewError(groupBaseCode+17, "get group details failed  "+groupName)
	ErrGroupUpdateGroupMemberRemark = errcode.NewError(groupBaseCode+18, "UpdateMemberRemark error failed  "+groupName)
	ErrGroupHandoverFailed          = errcode.NewError(groupBaseCode+19, "group handover  error failed  "+groupName)
	ErrGroupAssignAdmin             = errcode.NewError(groupBaseCode+20, "Group AssignAdmin  error failed  "+groupName)
	ErrGroupBanSpeakFailed          = errcode.NewError(groupBaseCode+21, "Group ban speak failed failed  "+groupName)
	ErrGroupMuteFailed              = errcode.NewError(groupBaseCode+22, "Group mute failed  "+groupName)
	ErrGroupApplyCreateFailed       = errcode.NewError(groupBaseCode+23, "Group apply create failed  "+groupName)
	ErrGroupApplyNotFound           = errcode.NewError(groupBaseCode+24, "Group apply not found  "+groupName)
	ErrGroupApplyAlreadyHandler     = errcode.NewError(groupBaseCode+25, "Group apply already handler  "+groupName)
	ErrGroupApplyUpdate             = errcode.NewError(groupBaseCode+26, "Group apply update err  "+groupName)
	ErrGroupNoticeDeleteFailed      = errcode.NewError(groupBaseCode+27, "Group notice delete err  "+groupName)
	ErrAllGroupNotice               = errcode.NewError(groupBaseCode+28, "get group all notice err  "+groupName)
	ErrNoticeCreateOrUpdateFailed   = errcode.NewError(groupBaseCode+29, "create or update notice err  "+groupName)

	// error codes are globally unique, adding 1 to the previous error code
)
