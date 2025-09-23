package protocol

import "buffersnow.com/spiritonline/pkg/gp"

const (
	GPCMCommand_KeepAlive       = "ka"
	GPCMCommand_Login           = "login"
	GPCMCommand_Logout          = "logout"
	GPCMCommand_UpdateProfile   = "updatepro"
	GPCMCommand_UpdateStatus    = "status"
	GPCMCommand_AddBuddy        = "addbuddy"
	GPCMCommand_DeleteBuddy     = "delbuddy"
	GPCMCommand_AuthorizeFriend = "authadd"
	GPCMCommand_BuddyMessage    = "bm"
	GPCMCommand_GetProfile      = "getprofile"
)

const (
	GPSPCommand_KeepAlive  = "ka"
	GPSPCommand_OthersList = "otherslist"
	GPSPCommand_Search     = "search"
)

var (
	GPError_Unknown            = gp.GameSpyError{ErrorCode: 0x0000, Message: "An unknown error occurred", IsFatal: true}
	GPError_Parse              = gp.GameSpyError{ErrorCode: 0x0001, Message: "Failed to processing incoming request", IsFatal: true}
	GPError_NeedsAuthorization = gp.GameSpyError{ErrorCode: 0x0002, Message: "Please login to make this request", IsFatal: true}
)

var (
	GPLoginError_Unknown = gp.GameSpyError{ErrorCode: 0x0040, Message: "There was an unknown login error", IsFatal: true}
)
