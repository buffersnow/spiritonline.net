package log

var levelComponentColors = []string{
	"\033[38;5;112m", // Hex: #87d700 intlLevel_Info
	"\033[38;5;153m", // Hex: #afd7ff intlLevel_Action
	"\033[38;5;73m",  // Hex: #5fafaf intlLevel_Event
	"\033[38;5;221m", // Hex: #ffd75f intlLevel_Warning
	"\033[38;5;160m", // Hex: #d70000 intlLevel_Error
	"\033[38;5;88m",  // Hex: #870000 intlLevel_Fatal
}

var debugComponentColors = []string{
	"\033[38;5;62m",  // Hex: #5f5fd7 DEBUG_GENERIC  -> any
	"\033[38;5;208m", // Hex: #ff8700 DEBUG_TRAFFIC  -> TCP/UDP/HTTP Connections
	"\033[38;5;93m",  // Hex: #8700ff DEBUG_SERVICE  -> Microservice
	"\033[38;5;191m", // Hex: #d7ff5f DEBUG_DATABASE -> GORM SQL/CouchDB queries
	"\033[38;5;111m", // Hex: #87afff DEBUG_API      -> Extra API logging (eg. WebAPI)
	"\033[38;5;212m", // Hex: #ff87d7 DEBUG_ROUTER   -> Spirit Internal
}

const (
	intlLevel_Info int = iota
	intlLevel_Action
	intlLevel_Event
	intlLevel_Warning
	intlLevel_Error
	intlLevel_Fatal
)

const (
	DEBUG_GENERIC int = iota
	DEBUG_TRAFFIC
	DEBUG_SERVICE
	DEBUG_DATABASE
	DEBUG_API
	DEBUG_ROUTER
)
