package hscore

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/lekuruu/hexagon/common"
)

type ScoreSubmissionRequest struct {
	Replay      []byte
	ProcessList string
	ScoreData   string
	Password    string
	ClientData  string
}

func (req *ScoreSubmissionRequest) String() string {
	return fmt.Sprintf(
		"ScoreSubmissionRequest{Replay: %d bytes, ProcessList: '%s', ScoreData: '%s', Password: '%s', ClientData: %s}",
		len(req.Replay),
		req.ProcessList,
		req.ScoreData,
		req.Password,
		req.ClientData,
	)
}

func NewScoreSubmissionRequest(request *http.Request) (*ScoreSubmissionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	replay := GetMultipartFormFile(request, "r")
	password := GetMultipartFormValue(request, "p")
	ivBase64 := GetMultipartFormValue(request, "iv")
	processListBase64 := GetMultipartFormValue(request, "pl")
	scoreDataBase64 := GetMultipartFormValue(request, "s")
	clientDataBase64 := GetMultipartFormValue(request, "i")

	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return nil, err
	}

	processList, err := base64.StdEncoding.DecodeString(processListBase64)
	if err != nil {
		return nil, err
	}

	scoreData, err := base64.StdEncoding.DecodeString(scoreDataBase64)
	if err != nil {
		return nil, err
	}

	clientData, err := base64.StdEncoding.DecodeString(clientDataBase64)
	if err != nil {
		return nil, err
	}

	processListDecrypted, err := common.DecryptScoreData(iv, processList)
	if err != nil {
		return nil, err
	}

	scoreDataDecrypted, err := common.DecryptScoreData(iv, scoreData)
	if err != nil {
		return nil, err
	}

	clientDataDecryptedBase64, err := common.DecryptScoreData(iv, clientData)
	if err != nil {
		return nil, err
	}

	clientDataDecrypted, err := base64.StdEncoding.DecodeString(string(clientDataDecryptedBase64))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// TODO: Parse process list, score data & client data
	return &ScoreSubmissionRequest{
		Replay:      replay,
		Password:    password,
		ProcessList: string(processListDecrypted),
		ScoreData:   string(scoreDataDecrypted),
		ClientData:  string(clientDataDecrypted),
	}, nil
}

func ScoreSubmissionHandler(ctx *Context) {
	if ctx.Request.Method != "POST" {
		ctx.Response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse score submission request
	req, err := NewScoreSubmissionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Errorf("Error parsing score submission request: %v", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Server.Logger.Infof("Score submission request: %s", req.String())
	// TODO: Process data & write response
}
