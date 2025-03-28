package hscore

import (
	"fmt"
	"net/http"
	"strconv"
)

func ReplayDownloadHandler(ctx *Context) {
	request, err := NewReplayDownloadRequest(*ctx.Request)
	if err != nil {
		ctx.Server.Logger.Errorf("Failed to parse replay download request: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	_, success := AuthenticateUser(
		request.Username,
		request.Password,
		ctx.Server,
	)

	if !success {
		ctx.Server.Logger.Warningf("Failed to authenticate user '%s'", request.Username)
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := ctx.Server.State.Storage.GetReplayFile(request.ScoreId)
	if err != nil {
		ctx.Server.Logger.Errorf("Failed to get replay file: %s", err)
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(file)
}

func NewReplayDownloadRequest(request http.Request) (*ReplayDownloadRequest, error) {
	query := request.URL.Query()

	username := query.Get("u")
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	password := query.Get("p")
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}

	scoreId := query.Get("s")
	if scoreId == "" {
		return nil, fmt.Errorf("missing score id")
	}

	scoreIdInt, err := strconv.Atoi(scoreId)
	if err != nil {
		return nil, fmt.Errorf("invalid score id: %s", scoreId)
	}

	return &ReplayDownloadRequest{
		Username: username,
		Password: password,
		ScoreId:  scoreIdInt,
	}, nil
}