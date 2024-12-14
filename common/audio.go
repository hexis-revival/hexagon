package common

import (
	"io"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// ExtractAudioSnippet extracts a snippet from the audio byte slice, compresses it, and saves it into an MP3 file
func ExtractAudioSnippet(audioData []byte, offsetMs int, durationMs int, bitrate int, storage Storage) ([]byte, error) {
	beatmapAudio, err := storage.CreateTempFile()
	if err != nil {
		return nil, err
	}

	_, err = beatmapAudio.Write(audioData)
	if err != nil {
		return nil, err
	}

	inputArgs := ffmpeg.KwArgs{"ss": offsetMs}
	outputArgs := ffmpeg.KwArgs{"t": durationMs, "ab": bitrate}
	ffmpeg.LogCompiledCommand = false

	err = ffmpeg.Input(beatmapAudio.Name(), inputArgs).
		Output(beatmapAudio.Name()+".mp3", outputArgs).
		Run()

	if err != nil {
		return nil, err
	}

	file, err := os.Open(beatmapAudio.Name() + ".mp3")
	if err != nil {
		return nil, err
	}

	defer func() {
		file.Close()
		os.Remove(beatmapAudio.Name())
		os.Remove(beatmapAudio.Name() + ".mp3")
	}()

	return io.ReadAll(file)
}
