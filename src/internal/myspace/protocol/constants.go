package protocol

const (
	MSIMCommand_BuddyMessage        = "bm"
	MSIMCommand_Error               = "error"
	MSIMCommand_LoginChallenge      = "lc"
	MSIMCommand_Keepalive           = "ka"
	MSIMCommand_CallbackRequest     = "persist"
	MSIMCommand_CallbackResponse    = "persistr"
	MSIMCommand_LoginResponseLegacy = "login"
	MSIMCommand_LoginResponse       = "login2"
	MSIMCommand_Logout              = "logout"
	MSIMCommand_UpdateStatus        = "status"
	MSIMCommand_AddBuddy            = "addbuddy"
	MSIMCommand_DeleteBuddy         = "delbuddy"
	MSIMCommand_UpdateBlocklist     = "blocklist"
	MSIMCommand_GetUserInfo         = "getinfo"
	MSIMCommand_SetUserInfo         = "setinfo"
	MSIMCommand_WebChallenge        = "webchlg"
	MSIMCommand_SkypeInterOp        = "skype"
	MSIMCommand_NetworkTest         = "nettest"
)

const (
	MSIMCallback_GetContactList          = "1;0;1"
	MSIMCallback_LookupCIbyUserId        = "1;0;2"
	MSIMCallback_LookupIMInfoForSelf     = "1;1;4"
	MSIMCallback_LookupIMInfoByUserId    = "1;1;17"
	MSIMCallback_GetContactGroupsList    = "1;2;6"
	MSIMCallback_LookupSMInfoForSelf     = "1;4;3"
	MSIMCallback_LookupSMInfoByUserId    = "1;4;5"
	MSIMCallback_LookupSMInfoByMail      = "1;5;7"
	MSIMCallback_HyperlinkRequest        = "1;6;11"
	MSIMCallback_QueryNotificationsForSM = "1;7;18"
	MSIMCallback_WebChallenge            = "1;17;26"
	MSIMCallback_GetFavouriteSongFromSM  = "1;21;18"
	MSIMCallback_GetServerConfiguration  = "1;101;20"
	MSIMCallback_SetProfilePicture       = "2;8;13"
	MSIMCallback_SetUsername             = "2;9;14"
	MSIMCallback_AddAllFriendsFromSM     = "2;14;21"
	MSIMCallback_AddTopFriendsFromSM     = "2;15;22"
	MSIMCallback_UpdateContactInfo       = "514;0;9"
	MSIMCallback_ChangeUserPreferences   = "514;1;10"
	MSIMCallback_SetProfilePictureLegacy = "514;8;13"
	MSIMCallback_InviteToIM              = "514;16;25"

	// 2;2;6 - Unknown
	// 2;2;16 - Unknown
	// 515;0;8 - Something about buddies deletion
)
