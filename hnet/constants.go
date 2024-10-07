package hnet

const (
	CLIENT_LOGIN               uint32 = 1
	CLIENT_CHANGE_STATUS       uint32 = 4
	CLIENT_REQUEST_STATS       uint32 = 11
	CLIENT_FOLLOW_USER         uint32 = 22
	CLIENT_UNFOLLOW_USER       uint32 = 23
	CLIENT_LEADERBOARD_REQUEST uint32 = 25
)

const (
	SERVER_LOGIN_RESPONSE uint32 = 2
	SERVER_USER_STATS     uint32 = 5
	SERVER_USER_INFO      uint32 = 6
	SERVER_FRIENDS_LIST   uint32 = 8
)

const (
	ACTION_IDLE       uint32 = 1
	ACTION_AWAY       uint32 = 2
	ACTION_PLAYING    uint32 = 3
	ACTION_EDITING    uint32 = 4
	ACTION_MODDING    uint32 = 5
	ACTION_TESTING    uint32 = 6
	ACTION_SUBMITTING uint32 = 7
	ACTION_WATCHING   uint32 = 8
)
