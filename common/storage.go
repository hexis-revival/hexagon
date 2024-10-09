package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Storage interface {
	// Base
	Save(key string, bucket string, data []byte) error
	Read(key string, bucket string) ([]byte, error)
	Remove(key string, bucket string) error
	Download(url string, key string, bucket string) error

	// Avatars
	GetAvatar(userId int) ([]byte, error)
	SaveAvatar(userId int, data []byte) error
	DefaultAvatar() ([]byte, error)
	EnsureDefaultAvatar() error
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

func (storage *FileStorage) GetAvatar(userId int) ([]byte, error) {
	avatar, err := storage.Read(fmt.Sprintf("%d", userId), "avatars")
	if err != nil {
		return storage.DefaultAvatar()
	}
	return avatar, nil
}

func (storage *FileStorage) SaveAvatar(userId int, data []byte) error {
	return storage.Save(string(userId), "avatars", data)
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
		return errors.New("failed to download default avatar")
	}

	return nil
}
