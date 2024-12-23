package hscore

import (
	"encoding/hex"

	"github.com/hexis-revival/hexagon/common"
)

func AuthenticateUser(username string, password string, server *ScoreServer) (*common.User, bool) {
	userObject, err := common.FetchUserByNameCaseInsensitive(
		username,
		server.State,
		"Stats",
	)

	if err != nil {
		server.Logger.Warningf("[Beatmap Submission] User '%s' not found", username)
		return nil, false
	}

	decodedPassword, err := hex.DecodeString(password)

	if err != nil {
		server.Logger.Warningf("[Beatmap Submission] Password decoding error: %s", err)
		return nil, false
	}

	isCorrect := common.CheckPasswordHashed(
		decodedPassword,
		userObject.Password,
	)

	if !isCorrect {
		server.Logger.Warningf("[Beatmap Submission] Incorrect password for '%s'", username)
		return nil, false
	}

	if !userObject.Activated {
		server.Logger.Warningf("[Beatmap Submission] Account not activated for '%s'", username)
		return nil, false
	}

	if userObject.Restricted {
		server.Logger.Warningf("[Beatmap Submission] Account restricted for '%s'", username)
		return nil, false
	}

	return userObject, true
}