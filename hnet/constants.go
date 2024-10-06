package hnet

const (
	CLIENT_LOGIN         uint32 = 1
	CLIENT_CHANGE_STATUS uint32 = 4
	CLIENT_REQUEST_STATS uint32 = 11
)

const (
	SERVER_LOGIN_RESPONSE uint32 = 2
	SERVER_USER_PRESENCE  uint32 = 6
)

const (
	ACTION_IDLE    uint32 = 1
	ACTION_PLAYING uint32 = 3
	ACTION_EDITING uint32 = 5
	ACTION_TESTING uint32 = 6
)
