package common

import (
	"encoding/binary"
	"io"
	"os"
	"testing"
)

func TestReplayParsing(t *testing.T) {
	file, err := os.Open("./replays_test.hxrp")
	if err != nil {
		t.Error(err)
		return
	}

	replayData, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
		return
	}

	stream := NewIOStream(replayData, binary.BigEndian)
	replay, err := ReadFullReplay(stream)
	if err != nil {
		t.Error(err)
		return
	}

	if replay.Header == nil {
		t.Error("failed to read replay header")
	}

	if len(replay.Frames) == 0 {
		t.Error("failed to read replay frames")
	}

	t.Cleanup(func() {
		file.Close()
	})
}
