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

	// Deserialize replay from file
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

	// Re-serialize replay
	stream = NewIOStream([]byte{}, binary.BigEndian)
	replay.Serialize(stream)

	if len(stream.Get()) == 0 {
		t.Error("failed to serialize replay")
	}

	// Deserialize replay from serialized data
	stream.Seek(0)
	replaySerialized, err := ReadFullReplay(stream)

	if err != nil {
		t.Error(err)
		return
	}

	if replaySerialized.Header == nil {
		t.Error("failed to read replay header")
	}

	if len(replaySerialized.Frames) == 0 {
		t.Error("failed to read replay frames")
	}
}
