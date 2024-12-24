package hnet

const (
	CLIENT_LOGIN               uint32 = 1
	CLIENT_LOGIN_RECONNECT     uint32 = 2
	CLIENT_CHANGE_STATUS       uint32 = 4
	CLIENT_REQUEST_STATS       uint32 = 11
	CLIENT_START_SPECTATING    uint32 = 14
	CLIENT_STOP_SPECTATING     uint32 = 15
	CLIENT_SPECTATE_HAS_MAP    uint32 = 16
	CLIENT_SPECTATE_FRAMES     uint32 = 18
	CLIENT_RELATIONSHIP_ADD    uint32 = 22
	CLIENT_RELATIONSHIP_REMOVE uint32 = 23
	CLIENT_LEADERBOARD_REQUEST uint32 = 25
	CLIENT_STATS_REFRESH       uint32 = 26
)

const (
	SERVER_LOGIN_RESPONSE         uint32 = 2
	SERVER_LOGIN_REVOKED          uint32 = 3
	SERVER_USER_STATS             uint32 = 5
	SERVER_USER_INFO              uint32 = 6
	SERVER_USER_QUIT              uint32 = 7
	SERVER_FRIENDS_LIST           uint32 = 8
	SERVER_SPECTATE_HAS_MAP       uint32 = 16
	SERVER_SPECTATE_STATUS_UPDATE uint32 = 17
	SERVER_SPECTATE_FRAMES        uint32 = 18
	SERVER_START_SPECTATING       uint32 = 19
	SERVER_STOP_SPECTATING        uint32 = 20
	SERVER_LEADERBOARD_RESPONSE   uint32 = 25
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

const (
	BEATMAP_STATUS_UNKNOWN       uint8 = 0
	BEATMAP_STATUS_NOT_SUBMITTED uint8 = 1
	BEATMAP_STATUS_PENDING       uint8 = 2
	BEATMAP_STATUS_RANKED        uint8 = 3
	BEATMAP_STATUS_APPROVED      uint8 = 4
)
