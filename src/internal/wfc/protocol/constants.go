package protocol

//% returncd codes
//% 200xx         = Hints
//% 201xx-209xx   = Error
//% 20xxxx and up = Poweroff (ie >999 ReCD)
//$ https://forum.wii-homebrew.com/index.php/Thread/51738-Collecting-Error-Codes/?pageNo=1#68adaac690392_0

const (
	ReCD_Login             = 001
	ReCD_AccountCreate     = 002
	ReCD_ServiceLocate     = 007
	ReCD_Maintenance       = 101  // unused, use *.AVAILABLE.GS.NINTENDOWIFI.NET bit 2
	ReCD_BannedFromWFC     = 102  // used by Wiimmfi and WiiLink, AltWFC uses >999
	ReCD_BrokenConID       = 103  // Invalid Base32, UserID and MAC don't match
	ReCD_ConIDInUse        = 104  // CID already in use on another console
	ReCD_MissingConID      = 105  // CID doesn't exist on server
	ReCD_InvalidPassword   = 105  // the DS uses this also for Invalid Password
	ReCD_TooManyUsers      = 106  // Too many user IDs created on console
	ReCD_UnsupportedGame   = 107  // Unsupported GameID
	ReCD_ConIDWasDeleted   = 108  // CID/UserID has been deleted
	ReCD_InvalidAction     = 109  // Wiimmfi uses this as Invalid GameID, WiiLink as Invalid Action
	ReCD_ServerShutdown    = 110  // Sent by nintendowifi.net after services went offline forever
	ReCD_ConIDAbuse        = 115  // Ie. this CID was seen too many times (eg. Dolphin w/o real NAND)
	ReCD_ServerUnavailable = 213  // Custom Error - Indicates that server is unavailable at the moment
	ReCD_UnknownConsole    = 343  // Custom Error - Indicates that the console hasn't been whitelisted yet
	ReCD_ConsolePending    = 365  // Custom Error - Indicates that the console is pending whitelist approval
	ReCD_PowerOffMessage   = 9999 // On Wii (and DS? needs testing) this will display a poweroff message
)
