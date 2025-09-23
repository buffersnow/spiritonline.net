package protocol

import "buffersnow.com/spiritonline/pkg/gp"

const (
	GPCommand_KeepAlive = "ka"
	GPCommand_Error     = "error"
)

const (
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
	GPSPCommand_OthersList = "otherslist"
	GPSPCommand_Search     = "search"
)

var (
	GPError_UnknownError       = &gp.GameSpyError{ErrorCode: 0x0000, Message: "An unknown error has occurred", IsFatal: true}
	GPError_Parse              = &gp.GameSpyError{ErrorCode: 0x0001, Message: "Failed to processing incoming request", IsFatal: true}
	GPError_UnknownCommand     = &gp.GameSpyError{ErrorCode: 0x0002, Message: "Unknown command sent to GSP backend", IsFatal: true}
	GPError_NeedsAuthorization = &gp.GameSpyError{ErrorCode: 0x0003, Message: "Please login to make this request", IsFatal: true}
)

var (
	GPLoginError_Unknown = &gp.GameSpyError{ErrorCode: 0x0040, Message: "There was an unknown login error", IsFatal: true}
)
