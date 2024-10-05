package common

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type ReplayFrame struct {
	Time        uint32
	MouseX      float64
	MouseY      float64
	ButtonState uint32
}

func (frame *ReplayFrame) String() string {
	return fmt.Sprintf(
		"ReplayFrame{Time: %d, MouseX: %f, MouseY: %f, ButtonState: %d}",
		frame.Time, frame.MouseX, frame.MouseY, frame.ButtonState,
	)
}

func (frame *ReplayFrame) Serialize(stream *IOStream) {
	stream.WriteU32(frame.Time)
	stream.WriteF64(frame.MouseX)
	stream.WriteF64(frame.MouseY)
	stream.WriteU32(frame.ButtonState)
}

func readReplayFrame(stream *IOStream) ReplayFrame {
	frame := ReplayFrame{}
	frame.Time = stream.ReadU32()
	frame.MouseX = stream.ReadF64()
	frame.MouseY = stream.ReadF64()
	frame.ButtonState = stream.ReadU32()
	return frame
}

type ReplayData struct {
	Frames []ReplayFrame
}

func (replayData *ReplayData) String() string {
	return fmt.Sprintf(
		"ReplayData{%d frames}",
		len(replayData.Frames),
	)
}

func (replayData *ReplayData) Serialize() []byte {
	stream := NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU32(uint32(len(replayData.Frames)))

	for _, frame := range replayData.Frames {
		frame.Serialize(stream)
	}

	decompressed := stream.Get()
	compressed := bytes.NewBuffer([]byte{})

	zlibWriter := zlib.NewWriter(compressed)
	zlibWriter.Write(decompressed)
	zlibWriter.Close()

	stream = NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU32(uint32(len(compressed.Bytes())))
	stream.Write(compressed.Bytes())
	return stream.Get()
}

func ReadCompressedReplay(replay []byte) (*ReplayData, error) {
	defer handlePanic()

	if len(replay) < 4 {
		return nil, fmt.Errorf("replay is too short")
	}

	stream := NewIOStream(replay, binary.BigEndian)
	replaySize := stream.ReadU32()
	compressedReplayData := stream.Read(int(replaySize))

	reader := bytes.NewReader(compressedReplayData)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}

	replayData, err := io.ReadAll(zlibReader)
	if err != nil {
		return nil, err
	}

	if len(replayData) < 4 {
		return nil, fmt.Errorf("replay data is too short")
	}

	stream = NewIOStream(replayData, binary.BigEndian)
	frameAmount := stream.ReadU32()
	frames := make([]ReplayFrame, frameAmount)

	// One frame is 24 bytes
	expectedSize := 24 * frameAmount

	// Check if we have enough data for all frames
	if stream.Available() < int(expectedSize) {
		return nil, fmt.Errorf(
			"not enough data for %d frames, got %d bytes",
			frameAmount, stream.Available(),
		)
	}

	for i := uint32(0); i < frameAmount; i++ {
		frames[i] = readReplayFrame(stream)
	}

	return &ReplayData{Frames: frames}, nil
}
