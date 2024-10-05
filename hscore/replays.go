package hscore

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lekuruu/hexagon/common"
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

func (frame *ReplayFrame) Serialize(stream *common.IOStream) {
	stream.WriteU32(frame.Time)
	stream.WriteF64(frame.MouseX)
	stream.WriteF64(frame.MouseY)
	stream.WriteU32(frame.ButtonState)
}

func ReadReplayFrame(stream *common.IOStream) ReplayFrame {
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
	stream := common.NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU32(uint32(len(replayData.Frames)))

	for _, frame := range replayData.Frames {
		frame.Serialize(stream)
	}

	decompressed := stream.Get()
	compressed := bytes.NewBuffer([]byte{})

	zlibWriter := zlib.NewWriter(compressed)
	zlibWriter.Write(decompressed)
	zlibWriter.Close()

	stream = common.NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU32(uint32(len(compressed.Bytes())))
	stream.Write(compressed.Bytes())
	return stream.Get()
}

func ReadCompressedReplay(replay []byte) (*ReplayData, error) {
	stream := common.NewIOStream(replay, binary.BigEndian)
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

	stream = common.NewIOStream(replayData, binary.BigEndian)
	frameAmount := stream.ReadU32()
	frames := make([]ReplayFrame, frameAmount)

	for i := uint32(0); i < frameAmount; i++ {
		frames[i] = ReadReplayFrame(stream)
	}

	return &ReplayData{Frames: frames}, nil
}
