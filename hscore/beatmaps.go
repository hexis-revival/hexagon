package hscore

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hexis-revival/hbxml"
	"github.com/hexis-revival/hexagon/common"
)

const (
	BssSuccess = 0
	// The username or password you entered is incorrect.
	BssAuthenticationError = 1
	// The beatmap you are trying to submit is no longer available.
	BssNotAvailable = 2
	// The beatmap you are trying to submit has already been ranked.
	BssAlreadyRanked = 3
	// The beatmap you are trying to submit is not owned by you.
	BssInvalidOwner = 4
)

func FormatBeatmapsetName(beatmapset *common.Beatmapset) string {
	return fmt.Sprintf(
		"%d %s - %s (%s)",
		beatmapset.Id,
		beatmapset.Artist,
		beatmapset.Title,
		beatmapset.Creator.Name,
	)
}

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

func CreateBeatmapset(beatmapIds []int, user *common.User, server *ScoreServer) (*common.Beatmapset, error) {
	beatmapset := &common.Beatmapset{
		Title:     "",
		Artist:    "",
		Source:    "",
		CreatorId: user.Id,
	}

	err := common.CreateBeatmapset(beatmapset, server.State)
	if err != nil {
		return nil, err
	}

	beatmaps := make([]common.Beatmap, 0, len(beatmapIds))

	for i := 0; i < len(beatmapIds); i++ {
		beatmap := common.Beatmap{SetId: beatmapset.Id, CreatorId: user.Id}
		beatmaps = append(beatmaps, beatmap)
	}

	err = common.CreateBeatmaps(beatmaps, server.State)
	if err != nil {
		return nil, err
	}

	beatmapset.Beatmaps = beatmaps
	return beatmapset, nil
}

func ValidateBeatmapset(beatmapset *common.Beatmapset, user *common.User, server *ScoreServer) int {
	if beatmapset.CreatorId != user.Id {
		server.Logger.Warningf("[Beatmap Submission] Invalid owner for '%s'", user.Name)
		return BssInvalidOwner
	}

	if beatmapset.Status > common.StatusPending {
		server.Logger.Warningf("[Beatmap Submission] Beatmapset already ranked (%d)", beatmapset.Id)
		return BssAlreadyRanked
	}

	if beatmapset.AvailabilityStatus != common.BeatmapHasDownload {
		server.Logger.Warningf("[Beatmap Submission] Beatmapset not available (%d)", beatmapset.Id)
		return BssNotAvailable
	}

	return BssSuccess
}

func UpdateBeatmapIds(beatmapset *common.Beatmapset, beatmapIds []int, server *ScoreServer) (_ []int, err error) {
	if len(beatmapIds) <= 0 {
		return nil, fmt.Errorf("beatmapIds list is empty")
	}

	currentBeatmapIds := make(map[int]bool, len(beatmapset.Beatmaps))
	beatmapIdMap := make(map[int]bool, len(beatmapIds))

	for _, beatmap := range beatmapset.Beatmaps {
		currentBeatmapIds[beatmap.Id] = true
	}

	for _, beatmapId := range beatmapIds {
		beatmapIdMap[beatmapId] = true
	}

	if len(beatmapIds) < len(beatmapset.Beatmaps) {
		// Ensure every beatmap id is inside current beatmap ids
		for _, beatmapId := range beatmapIds {
			if !currentBeatmapIds[beatmapId] {
				return nil, fmt.Errorf("Beatmap '%d' not found", beatmapId)
			}
		}

		// Remove unused beatmaps
		for _, beatmap := range beatmapset.Beatmaps {
			if !beatmapIdMap[beatmap.Id] {
				continue
			}

			err := common.RemoveBeatmap(&beatmap, server.State)
			if err != nil {
				return nil, err
			}
		}

		server.Logger.Debugf(
			"Removed %d beatmaps from beatmapset",
			len(currentBeatmapIds)-len(beatmapIds),
		)
	}

	// Calculate how many beatmaps we need to create
	requiredMaps := max(0, len(beatmapIds)-len(currentBeatmapIds))
	newBeatmaps := make([]common.Beatmap, 0, requiredMaps)

	// Create new beatmaps
	for i := 0; i < requiredMaps; i++ {
		beatmap := common.Beatmap{SetId: beatmapset.Id, CreatorId: beatmapset.CreatorId}
		newBeatmaps = append(newBeatmaps, beatmap)

		err := common.CreateBeatmap(&beatmap, server.State)
		if err != nil {
			return nil, err
		}
	}

	server.Logger.Debugf(
		"Added %d beatmaps to beatmapset",
		requiredMaps,
	)

	// Append new beatmaps to the beatmapset & return new beatmap ids
	beatmapset.Beatmaps, err = common.FetchBeatmapsBySetId(
		beatmapset.Id,
		server.State,
	)

	if err != nil {
		return nil, err
	}

	beatmapIds = make([]int, 0, len(beatmapset.Beatmaps))

	for _, beatmap := range beatmapset.Beatmaps {
		beatmapIds = append(beatmapIds, beatmap.Id)
	}

	return beatmapIds, nil
}

func RemoveInactiveBeatmaps(user *common.User, server *ScoreServer) error {
	beatmapsets, err := common.FetchBeatmapsetsByStatus(
		user.Id,
		common.StatusNotSubmitted,
		server.State,
	)

	if err != nil {
		return err
	}

	server.Logger.Debugf(
		"Found %d inactive beatmapsets for '%s'",
		len(beatmapsets), user.Name,
	)
	// TODO: Remove beatmap assets from storage

	for _, beatmapset := range beatmapsets {
		err := common.RemoveBeatmapsBySetId(beatmapset.Id, server.State)
		if err != nil {
			return err
		}

		err = common.RemoveBeatmapset(&beatmapset, server.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemainingBeatmapUploads(user *common.User, server *ScoreServer) (int, error) {
	unrankedBeatmaps, err := common.FetchBeatmapsetUnrankedCountByCreatorId(
		user.Id,
		server.State,
	)

	if err != nil {
		return 0, err
	}

	rankedBeatmaps, err := common.FetchBeatmapsetRankedCountByCreatorId(
		user.Id,
		server.State,
	)

	if err != nil {
		return 0, err
	}

	// Users can upload up to 6 pending maps plus
	// 1 per ranked map, up to a maximum of 10.
	return (6 - unrankedBeatmaps) + min(rankedBeatmaps, 10), nil
}

func ProcessUploadPackage(request *BeatmapUploadRequest, server *ScoreServer) (map[string][]byte, map[string]*hbxml.Beatmap, error) {
	files := make(map[string][]byte, 0)
	beatmaps := make(map[string]*hbxml.Beatmap, 0)

	for _, file := range request.Package.File {
		fileHandle, err := request.Package.Open(file.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open file: %s", err)
		}

		files[file.Name], err = io.ReadAll(fileHandle)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read file: %s", err)
		}

		if !strings.HasSuffix(file.Name, ".hbxml") {
			continue
		}

		reader := bytes.NewReader(files[file.Name])
		beatmaps[file.Name], err = hbxml.NewBeatmap(reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse beatmap: %s", err)
		}
	}

	drainLength := GetMaximumDrainLength(beatmaps)
	maximumSize := CalculatePackageSizeLimit(drainLength)
	packageSize := 0

	for _, file := range files {
		packageSize += len(file)
	}

	if packageSize > int(maximumSize) {
		return nil, nil, errors.New("package size limit exceeded")
	}

	return files, beatmaps, nil
}

func ResolveBeatmaps(beatmapObjects map[string]*hbxml.Beatmap, beatmaps []common.Beatmap) (map[string]*common.Beatmap, error) {
	if len(beatmapObjects) != len(beatmaps) {
		return nil, errors.New("beatmap count mismatch")
	}

	beatmapMap := make(map[string]*common.Beatmap, 0)
	foundBeatmaps := make(map[string]bool, len(beatmaps))

	// Find existing beatmaps and map them
	for filename, _ := range beatmapObjects {
		for _, beatmap := range beatmaps {
			if beatmap.Filename == filename {
				beatmapMap[filename] = &beatmap
				foundBeatmaps[filename] = true
				break
			}
		}
	}

	missingBeatmaps := make([]string, 0)

	for filename, _ := range beatmapObjects {
		if foundBeatmaps[filename] {
			continue
		}

		missingBeatmaps = append(missingBeatmaps, filename)
	}

	// Assign missing beatmaps
	for _, beatmap := range beatmaps {
		if foundBeatmaps[beatmap.Filename] {
			continue
		}

		if len(missingBeatmaps) <= 0 {
			return nil, errors.New("missing beatmaps")
		}

		beatmapMap[missingBeatmaps[0]] = &beatmap
		missingBeatmaps = missingBeatmaps[1:]
	}

	return beatmapMap, nil
}

func CalculatePackageSizeLimit(beatmapLength int) float64 {
	// The file size limit is 10MB plus an additional 10MB for
	// every minute of beatmap length, and it caps at 100MB.
	return math.Min(
		float64(10_000_000+(10_000_000*(beatmapLength/60))),
		100_000_000,
	)
}

func GetMaximumDrainLength(beatmapObjects map[string]*hbxml.Beatmap) int {
	maxDrainLength := 0

	for _, beatmap := range beatmapObjects {
		maxDrainLength = max(maxDrainLength, int(beatmap.DrainLength()))
	}

	return maxDrainLength
}

func HasExtension(filename string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}

	return false
}

func UpdateBeatmapsetMetadata(beatmapset *common.Beatmapset, metadata hbxml.Meta, server *ScoreServer) error {
	beatmapset.Title = metadata.Title
	beatmapset.Artist = metadata.Artist
	beatmapset.Source = metadata.Source
	beatmapset.Tags = metadata.Tags
	beatmapset.LastUpdated = time.Now()
	beatmapset.Status = common.StatusPending

	if beatmapset.Tags == nil {
		beatmapset.Tags = []string{}
	}

	return common.UpdateBeatmapset(beatmapset, server.State)
}

func UpdateBeatmapMetadata(beatmap *common.Beatmap, beatmapObject *hbxml.Beatmap, file []byte, filename string, server *ScoreServer) error {
	beatmapChecksumBytes := md5.Sum(file)
	beatmap.Filename = filename
	beatmap.Status = common.StatusPending
	beatmap.Checksum = hex.EncodeToString(beatmapChecksumBytes[:])
	beatmap.Version = beatmapObject.Meta.Version
	beatmap.TotalLength = int(beatmapObject.TotalLength())
	beatmap.DrainLength = int(beatmapObject.DrainLength())
	beatmap.TotalCircles = beatmapObject.TotalCircles()
	beatmap.TotalSliders = beatmapObject.TotalSliders()
	beatmap.TotalSpinners = beatmapObject.TotalSpinners()
	beatmap.TotalHolds = beatmapObject.TotalHolds()
	beatmap.MedianBpm = beatmapObject.MedianBPM()
	beatmap.HighestBpm = beatmapObject.HighestBPM()
	beatmap.LowestBpm = beatmapObject.LowestBPM()
	beatmap.CS = beatmapObject.Difficulty.CircleSize
	beatmap.HP = beatmapObject.Difficulty.HPDrainRate
	beatmap.OD = beatmapObject.Difficulty.OverallDifficulty
	beatmap.AR = beatmapObject.Difficulty.ApproachRate
	beatmap.LastUpdated = time.Now()
	// TODO: Update MaxCombo & Star Rating
	return common.UpdateBeatmap(beatmap, server.State)
}

func UploadBeatmapPackage(setId int, files map[string][]byte, server *ScoreServer) error {
	allowedFileExtensions := []string{
		".hbxml", ".hxz", ".png", ".jpg", ".jpeg", ".mp3",
		".ogg", ".wav", ".flac", ".wmv", ".flv", ".mp4",
		".avi", ".mkv", ".webm", ".txt", ".ini",
	}

	buffer := bytes.Buffer{}
	zipWriter := zip.NewWriter(&buffer)

	for filename, file := range files {
		if !HasExtension(filename, allowedFileExtensions) {
			continue
		}

		fileWriter, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}

		_, err = fileWriter.Write(file)
		if err != nil {
			return err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return err
	}

	// TODO: Validate package files
	// TODO: Limit package size
	return server.State.Storage.SaveBeatmapPackage(
		setId,
		buffer.Bytes(),
	)
}

func UploadBeatmapFiles(beatmapFiles map[int][]byte, server *ScoreServer) error {
	for beatmapId, file := range beatmapFiles {
		err := server.State.Storage.SaveBeatmapFile(beatmapId, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func UploadAudioPreview(setId int, files map[string][]byte, general hbxml.General, server *ScoreServer) error {
	offset := general.PreviewOffset / 1000
	audioFilename := general.AudioFilename

	if offset < 0 {
		return errors.New("invalid preview offset")
	}

	audio, ok := files[audioFilename]
	if !ok {
		return errors.New("audio file not found")
	}

	audioSnippet, err := common.ExtractAudioSnippet(
		audio,
		offset,
		10,
		64,
		server.State.Storage,
	)

	if err != nil {
		return err
	}

	return server.State.Storage.SaveBeatmapPreview(setId, audioSnippet)
}

func UploadBeatmapThumbnail(setId int, files map[string][]byte, events hbxml.Events, server *ScoreServer) error {
	if events.Backgrounds == nil {
		// No background found
		return nil
	}

	background := events.Backgrounds[0]
	backgroundFilename := background.Filename

	imageData, ok := files[backgroundFilename]
	if !ok {
		return errors.New("background file not found")
	}

	generator := common.NewImageGenerator(common.ImageGenerator{})
	generator.Scaler = "CatmullRom"
	generator.Width = 160
	generator.Height = 120

	image, err := generator.NewImageFromByteArray(imageData)
	if err != nil {
		return err
	}

	largeImage, err := generator.CreateThumbnail(image)
	if err != nil {
		return err
	}

	generator.Width = 80
	generator.Height = 60

	smallImage, err := generator.CreateThumbnail(image)
	if err != nil {
		return err
	}

	err = server.State.Storage.SaveBeatmapThumbnail(setId, largeImage, true)
	if err != nil {
		return err
	}

	err = server.State.Storage.SaveBeatmapThumbnail(setId, smallImage, false)
	if err != nil {
		return err
	}

	return nil
}

func BeatmapGenIdHandler(ctx *Context) {
	request, err := NewBeatmapSubmissionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	user, success := AuthenticateUser(
		request.Username,
		request.Password,
		ctx.Server,
	)

	response := &BeatmapSubmissionResponse{
		StatusCode: BssSuccess,
		SetId:      -1,
		BeatmapIds: []int{},
	}

	if !success {
		response.StatusCode = BssAuthenticationError
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	err = RemoveInactiveBeatmaps(user, ctx.Server)
	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to remove inactive beatmaps: %s", err)
	}

	beatmapset, err := common.FetchBeatmapsetById(
		request.SetId,
		ctx.Server.State,
		"Beatmaps",
	)

	if err != nil {
		if err.Error() != "record not found" {
			// Database failure
			ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset fetch error: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		remainingUploads, err := RemainingBeatmapUploads(user, ctx.Server)
		if err != nil {
			ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to fetch remaining uploads: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		if remainingUploads <= 0 {
			ctx.Server.Logger.Warningf("[Beatmap Submission] No remaining uploads for '%s'", user.Name)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		beatmapset, err = CreateBeatmapset(
			request.BeatmapIds,
			user,
			ctx.Server,
		)

		if err != nil {
			ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset creation error: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		for _, beatmap := range beatmapset.Beatmaps {
			response.BeatmapIds = append(response.BeatmapIds, beatmap.Id)
		}

		response.SetId = beatmapset.Id
		response.StatusCode = BssSuccess
		ctx.Response.Write([]byte(response.Write()))
		ctx.Server.Logger.Infof(
			"[Beatmap Submission] Beatmapset for '%s' created with %d beatmaps (%d)",
			user.Name,
			len(beatmapset.Beatmaps),
			beatmapset.Id,
		)
		return
	}

	response.StatusCode = ValidateBeatmapset(
		beatmapset,
		user,
		ctx.Server,
	)

	if response.StatusCode != BssSuccess {
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	response.BeatmapIds, err = UpdateBeatmapIds(
		beatmapset,
		request.BeatmapIds,
		ctx.Server,
	)

	if err != nil {
		ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset update error: %s", err)
		response.StatusCode = BssNotAvailable
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	response.SetId = beatmapset.Id
	ctx.Server.Logger.Infof(
		"[Beatmap Submission] Beatmapset for '%s' updated with %d beatmaps (%d)",
		user.Name,
		len(beatmapset.Beatmaps),
		beatmapset.Id,
	)
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapUploadHandler(ctx *Context) {
	request, err := NewBeatmapUploadRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Upload request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Upload request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)
	response := &BeatmapUploadResponse{Success: true}

	user, success := AuthenticateUser(
		request.Username,
		request.Password,
		ctx.Server,
	)

	if !success {
		response.Success = false
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	beatmapset, err := common.FetchBeatmapsetById(
		request.SetId,
		ctx.Server.State,
		"Beatmaps",
	)

	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Beatmapset fetch error: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	statusCode := ValidateBeatmapset(
		beatmapset,
		user,
		ctx.Server,
	)

	if statusCode != BssSuccess {
		response.Success = false
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	files, beatmapObjects, err := ProcessUploadPackage(
		request,
		ctx.Server,
	)

	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	beatmapMap, err := ResolveBeatmaps(
		beatmapObjects,
		beatmapset.Beatmaps,
	)

	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Got %d files in package.", len(files))

	var general hbxml.General
	var metadata hbxml.Meta
	var events hbxml.Events

	for filename, beatmap := range beatmapObjects {
		general = beatmap.General
		metadata = beatmap.Meta
		events = beatmap.Events

		err = UpdateBeatmapMetadata(
			beatmapMap[filename],
			beatmapObjects[filename],
			files[filename],
			filename,
			ctx.Server,
		)

		if err != nil {
			response.Success = false
			ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to update beatmap metadata: %s", err)
			ctx.Response.Write([]byte(response.Write()))
			return
		}
	}

	beatmapIdMap := make(map[int][]byte, 0)

	for _, beatmap := range beatmapset.Beatmaps {
		beatmapIdMap[beatmap.Id] = files[beatmap.Filename]
	}

	err = UpdateBeatmapsetMetadata(
		beatmapset,
		metadata,
		ctx.Server,
	)

	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to update beatmapset metadata: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	err = UploadBeatmapFiles(beatmapIdMap, ctx.Server)
	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to upload beatmap files: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	err = UploadBeatmapPackage(beatmapset.Id, files, ctx.Server)
	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to upload beatmap package: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	err = UploadAudioPreview(beatmapset.Id, files, general, ctx.Server)
	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to upload audio preview: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	err = UploadBeatmapThumbnail(beatmapset.Id, files, events, ctx.Server)
	if err != nil {
		response.Success = false
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to upload beatmap thumbnail: %s", err)
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	ctx.Server.Logger.Infof(
		"[Beatmap Submission] Beatmapset '%s' updated by '%s'",
		FormatBeatmapsetName(beatmapset),
		user.Name,
	)

	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapGenTopicHandler(ctx *Context) {
	request, err := NewBeatmapDescriptionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Description request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Description request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	response := &BeatmapDescriptionResponse{
		TopicId: 1,
		Content: "The quick brown fox jumps over the lazy dog.",
	}

	// TODO: Implement beatmap description logic
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapPostHandler(ctx *Context) {
	request, err := NewBeatmapPostRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Post request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Post request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)
	// TODO: Implement beatmap post logic
}

func NewBeatmapSubmissionRequest(request *http.Request) (*BeatmapSubmissionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	setId := GetMultipartFormValue(request, "s")
	beatmapIdsString := GetMultipartFormValue(request, "v")
	clientVersion := GetMultipartFormValue(request, "x")

	collection := common.NewErrorCollection()

	setIdInt, err := strconv.Atoi(setId)
	collection.Add(err)

	clientVersionInt, err := strconv.Atoi(clientVersion)
	collection.Add(err)

	beatmapIdsList := strings.Split(beatmapIdsString, ":")
	beatmapIds := make([]int, 0, len(beatmapIdsList))

	for _, beatmapId := range beatmapIdsList {
		beatmapIdInt, err := strconv.Atoi(beatmapId)
		beatmapIds = append(beatmapIds, beatmapIdInt)
		collection.Add(err)
	}

	if collection.HasErrors() {
		return nil, collection.Pop(0)
	}

	return &BeatmapSubmissionRequest{
		Username:      username,
		Password:      password,
		BeatmapIds:    beatmapIds,
		SetId:         setIdInt,
		ClientVersion: clientVersionInt,
	}, nil
}

func NewBeatmapUploadRequest(request *http.Request) (*BeatmapUploadRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")
	setId := GetMultipartFormValue(request, "s")

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, err
	}

	zip, err := GetMultipartZipFile(request, "d")
	if err != nil {
		return nil, err
	}

	return &BeatmapUploadRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		SetId:         setIdInt,
		Package:       zip,
	}, nil
}

func NewBeatmapDescriptionRequest(request *http.Request) (*BeatmapDescriptionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")
	setId := GetMultipartFormValue(request, "s")

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, err
	}

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	return &BeatmapDescriptionRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		SetId:         setIdInt,
	}, nil
}

func NewBeatmapPostRequest(request *http.Request) (*BeatmapPostRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")
	setId := GetMultipartFormValue(request, "s")
	content := GetMultipartFormValue(request, "t")

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, err
	}

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	return &BeatmapPostRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		SetId:         setIdInt,
		Content:       content,
	}, nil
}
