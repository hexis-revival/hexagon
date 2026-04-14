package common

import (
	"encoding/binary"
	"io"
	"os"
	"reflect"
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

	if replay.Mods == nil {
		t.Fatal("expected parsed replay mods")
	}

	if replay.Checksum() != replay.ScoreChecksum {
		t.Fatalf(
			"replay checksum mismatch: got %s want %s",
			replay.Checksum(),
			replay.ScoreChecksum,
		)
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

	if replaySerialized.ScoreChecksum != replay.ScoreChecksum {
		t.Fatalf(
			"re-serialized replay checksum mismatch: got %s want %s",
			replaySerialized.ScoreChecksum,
			replay.ScoreChecksum,
		)
	}

	if len(replaySerialized.Frames) == 0 {
		t.Error("failed to read replay frames")
	}

	if !reflect.DeepEqual(replaySerialized, replay) {
		t.Fatal("re-serialized replay does not match original")
	}
}
