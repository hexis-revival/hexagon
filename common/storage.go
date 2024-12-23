package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type Storage interface {
	// Base
	Save(key string, bucket string, data []byte) error
	Read(key string, bucket string) ([]byte, error)
	Remove(key string, bucket string) error
	Download(url string, key string, bucket string) error
	CreateTempFile() (*os.File, error)

	// Replays
	GetReplayFile(replayId int) ([]byte, error)
	SaveReplayFile(replayId int, data []byte) error
	RemoveReplayFile(replayId int) error

	// Avatars
	GetAvatar(userId int) ([]byte, error)
	SaveAvatar(userId int, data []byte) error
	DefaultAvatar() ([]byte, error)
	EnsureDefaultAvatar() error

	// Beatmaps
	GetBeatmapFile(beatmapId int) ([]byte, error)
	GetBeatmapPackage(beatmapsetId int) ([]byte, error)
	GetBeatmapThumbnail(beatmapId int, large bool) ([]byte, error)
	GetBeatmapPreview(beatmapId int) ([]byte, error)
	SaveBeatmapFile(beatmapId int, data []byte) error
	SaveBeatmapPackage(beatmapsetId int, data []byte) error
	SaveBeatmapThumbnail(beatmapId int, data []byte, large bool) error
	SaveBeatmapPreview(beatmapId int, data []byte) error
	RemoveBeatmapFile(beatmapId int) error
	RemoveBeatmapPackage(beatmapsetId int) error
	RemoveBeatmapThumbnail(beatmapId int) error
	RemoveBeatmapPreview(beatmapId int) error
}

type FileStorage struct {
	dataPath string
}

func NewFileStorage(dataPath string) Storage {
	return &FileStorage{dataPath: dataPath}
}

func (storage *FileStorage) Read(key string, folder string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.ReadFile(path)
}

func (storage *FileStorage) Save(key string, folder string, data []byte) error {
	path := fmt.Sprintf("%s/%s", storage.dataPath, folder)
	err := os.MkdirAll(path, 0755)

	if err != nil {
		return err
	}

	os.WriteFile(fmt.Sprintf("%s/%s", path, key), data, os.ModePerm)
	return nil
}

func (storage *FileStorage) Remove(key string, folder string) error {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Remove(path)
}

func (storage *FileStorage) Download(url string, key string, folder string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download failed: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return storage.Save(key, folder, data)
}

func (storage *FileStorage) CreateTempFile() (*os.File, error) {
	return os.CreateTemp(storage.dataPath, "temp")
}

func (storage *FileStorage) GetReplayFile(scoreId int) ([]byte, error) {
	return storage.Read(fmt.Sprintf("%d", scoreId), "replays")
}

func (storage *FileStorage) SaveReplayFile(scoreId int, data []byte) error {
	return storage.Save(strconv.Itoa(scoreId), "replays", data)
}

func (storage *FileStorage) RemoveReplayFile(scoreId int) error {
	return storage.Remove(strconv.Itoa(scoreId), "replays")
}

func (storage *FileStorage) GetAvatar(userId int) ([]byte, error) {
	avatar, err := storage.Read(fmt.Sprintf("%d", userId), "avatars")
	if err != nil {
		return storage.DefaultAvatar()
	}
	return avatar, nil
}

func (storage *FileStorage) SaveAvatar(userId int, data []byte) error {
	return storage.Save(strconv.Itoa(userId), "avatars", data)
}

func (storage *FileStorage) DefaultAvatar() ([]byte, error) {
	return storage.Read("unknown", "avatars")
}

func (storage *FileStorage) EnsureDefaultAvatar() error {
	_, err := storage.DefaultAvatar()
	if err == nil {
		return nil
	}

	// Download the default avatar
	err = storage.Download(
		"https://raw.githubusercontent.com/hexis-revival/hexagon/refs/heads/main/.github/images/unknown.png",
		"unknown", "avatars",
	)

	if err != nil {
		return errors.Wrap(err, "failed to download default avatar")
	}

	return nil
}

func (storage *FileStorage) GetBeatmapFile(beatmapId int) ([]byte, error) {
	return storage.Read(fmt.Sprintf("%d", beatmapId), "beatmaps")
}

func (storage *FileStorage) GetBeatmapPackage(beatmapsetId int) ([]byte, error) {
	return storage.Read(fmt.Sprintf("%d", beatmapsetId), "packages")
}

func (storage *FileStorage) GetBeatmapThumbnail(beatmapId int, large bool) ([]byte, error) {
	return storage.Read(formatThumbnailName(beatmapId, large), "thumbnails")
}

func (storage *FileStorage) GetBeatmapPreview(beatmapId int) ([]byte, error) {
	return storage.Read(fmt.Sprintf("%d", beatmapId), "previews")
}

func (storage *FileStorage) SaveBeatmapFile(beatmapId int, data []byte) error {
	return storage.Save(strconv.Itoa(beatmapId), "beatmaps", data)
}

func (storage *FileStorage) SaveBeatmapPackage(beatmapsetId int, data []byte) error {
	return storage.Save(strconv.Itoa(beatmapsetId), "packages", data)
}

func (storage *FileStorage) SaveBeatmapThumbnail(beatmapId int, data []byte, large bool) error {
	return storage.Save(formatThumbnailName(beatmapId, large), "thumbnails", data)
}

func (storage *FileStorage) SaveBeatmapPreview(beatmapId int, data []byte) error {
	return storage.Save(strconv.Itoa(beatmapId), "previews", data)
}

func (storage *FileStorage) RemoveBeatmapFile(beatmapId int) error {
	return storage.Remove(strconv.Itoa(beatmapId), "beatmaps")
}

func (storage *FileStorage) RemoveBeatmapPackage(beatmapsetId int) error {
	return storage.Remove(strconv.Itoa(beatmapsetId), "packages")
}

func (storage *FileStorage) RemoveBeatmapThumbnail(beatmapId int) error {
	err := storage.Remove(formatThumbnailName(beatmapId, true), "thumbnails")
	if err != nil {
		return err
	}
	return storage.Remove(formatThumbnailName(beatmapId, false), "thumbnails")
}

func (storage *FileStorage) RemoveBeatmapPreview(beatmapId int) error {
	return storage.Remove(strconv.Itoa(beatmapId), "previews")
}

func formatThumbnailName(beatmapId int, large bool) string {
	return fmt.Sprintf("%d%s", beatmapId, getThumbnailSuffix(large))
}

func getThumbnailSuffix(large bool) string {
	if large {
		return "_large"
	}
	return "_small"
}
