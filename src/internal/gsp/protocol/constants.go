package protocol

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
