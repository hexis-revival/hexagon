package common

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

type ReplayData struct {
	Header *ReplayHeader
	Frames []*ReplayFrame
}

func (replayData *ReplayData) String() string {
	header := "no header"

	if replayData.Header != nil {
		header = replayData.Header.String()
	}

	return fmt.Sprintf(
		"ReplayData{%d frames, %s}",
		len(replayData.Frames), header,
	)
}

func (replayData *ReplayData) SerializeFrames(stream *IOStream) {
	replayStream := NewIOStream([]byte{}, binary.BigEndian)
	replayStream.WriteU32(uint32(len(replayData.Frames)))

	for _, frame := range replayData.Frames {
		frame.Serialize(replayStream)
	}

	decompressed := replayStream.Get()
	compressed := bytes.NewBuffer([]byte{})

	zlibWriter := zlib.NewWriter(compressed)
	zlibWriter.Write(decompressed)
	zlibWriter.Close()

	stream.WriteU32(uint32(len(compressed.Bytes())))
	stream.Write(compressed.Bytes())
}

func (replayData *ReplayData) Serialize(stream *IOStream) {
	replayData.Header.Serialize(stream)
	replayData.SerializeFrames(stream)
}

type ReplayHeader struct {
	ReplayVersion   uint8
	BeatmapChecksum string
	PlayerName      string
	ScoreChecksum   string
	Count300        uint32
	Count100        uint32
	Count50         uint32
	CountGeki       uint32
	CountGood       uint32
	CountMiss       uint32
	TotalScore      float64
	MaxCombo        uint32
	FullCombo       bool
	Time            time.Time
	ModsData        [9]byte
}

func (header *ReplayHeader) String() string {
	return fmt.Sprintf(
		"ReplayHeader{ReplayVersion: %d, BeatmapChecksum: %s, PlayerName: %s, ScoreChecksum: %s, Count300: %d, Count100: %d, Count50: %d, CountGeki: %d, CountGood: %d, CountMiss: %d, TotalScore: %f, MaxCombo: %d, FullCombo: %t, Time: %s}",
		header.ReplayVersion, header.BeatmapChecksum, header.PlayerName, header.ScoreChecksum, header.Count300, header.Count100, header.Count50, header.CountGeki, header.CountGood, header.CountMiss, header.TotalScore, header.MaxCombo, header.FullCombo, header.Time,
	)
}

func (header *ReplayHeader) Serialize(stream *IOStream) {
	stream.Write([]byte{0xf0, 0xa0, 0xc0, 0xe0})
	stream.WriteU32(0)
	stream.WriteU8(header.ReplayVersion)
	stream.WriteString(header.BeatmapChecksum)
	stream.WriteString(header.PlayerName)
	stream.WriteString(header.ScoreChecksum)
	stream.WriteU32(header.Count300)
	stream.WriteU32(header.Count100)
	stream.WriteU32(header.Count50)
	stream.WriteU32(header.CountGeki)
	stream.WriteU32(header.CountGood)
	stream.WriteU32(header.CountMiss)
	stream.WriteF64(header.TotalScore)
	stream.WriteU32(header.MaxCombo)
	stream.WriteBool(header.FullCombo)
	stream.WriteDateTime(header.Time)
	stream.Write(header.ModsData[:])
}

func (header *ReplayHeader) Accuracy() float64 {
	totalHits := header.Count300 + header.Count100 + header.Count50
	return float64(header.Count300*300+header.Count100*100+header.Count50*50) / float64(totalHits*300)
}

func (header *ReplayHeader) Grade() Grade {
	totalHits := header.Count300 + header.Count100 + header.Count50 + header.CountGood

	if totalHits == 0 {
		return GradeF
	}

	totalHitCount := float64(totalHits)
	accuracyRatio := float64(header.Count300) / totalHitCount

	if math.IsNaN(accuracyRatio) || accuracyRatio == 1.0 {
		// TODO: Check if hidden is enabled
		// if header.Mods.Hidden {
		// 	  return GradeXH
		// }
		return GradeSS
	}

	if accuracyRatio <= 0.8 && header.CountGood == 0 {
		if accuracyRatio > 0.6 {
			return GradeC
		}
		return GradeD
	}

	if accuracyRatio <= 0.9 {
		return GradeB
	}

	// Default case for remaining conditions
	return GradeA
}

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

func ReadReplayFrame(stream *IOStream) *ReplayFrame {
	frame := &ReplayFrame{}
	frame.Time = stream.ReadU32()
	frame.MouseX = stream.ReadF64()
	frame.MouseY = stream.ReadF64()
	frame.ButtonState = stream.ReadU32()
	return frame
}

func ReadReplayHeader(stream *IOStream) *ReplayHeader {
	header := &ReplayHeader{}
	magicBytes := stream.Read(4)

	if !bytes.Equal(magicBytes, []byte{0xf0, 0xa0, 0xc0, 0xe0}) {
		return nil
	}

	_ = stream.ReadI32() // Always zero-bytes, we can ignore them
	header.ReplayVersion = stream.ReadU8()
	header.BeatmapChecksum = stream.ReadString()
	header.PlayerName = stream.ReadString()
	header.ScoreChecksum = stream.ReadString()
	header.Count300 = stream.ReadU32()
	header.Count100 = stream.ReadU32()
	header.Count50 = stream.ReadU32()
	header.CountGeki = stream.ReadU32()
	header.CountGood = stream.ReadU32()
	header.CountMiss = stream.ReadU32()
	header.TotalScore = stream.ReadF64()
	header.MaxCombo = stream.ReadU32()
	header.FullCombo = stream.ReadBool()
	header.Time = stream.ReadDateTime()
	header.ModsData = [9]byte(stream.Read(9))
	return header
}

func ReadCompressedReplay(stream *IOStream) (*ReplayData, error) {
	defer handlePanic()

	if stream.Available() < 4 {
		return nil, fmt.Errorf("replay is too short")
	}

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
	frames := make([]*ReplayFrame, frameAmount)

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
		frames[i] = ReadReplayFrame(stream)
	}

	return &ReplayData{Frames: frames}, nil
}

func ReadFullReplay(stream *IOStream) (*ReplayData, error) {
	header := ReadReplayHeader(stream)

	if header == nil {
		return &ReplayData{}, fmt.Errorf("failed to read replay header")
	}

	replayData, err := ReadCompressedReplay(stream)

	if err != nil {
		return &ReplayData{}, err
	}

	replayData.Header = header
	return replayData, nil
}
