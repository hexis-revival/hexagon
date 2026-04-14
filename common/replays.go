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
	Mode            int8
	ReplayVersion   int
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
	TimeSpec        int8

	Frames []*ReplayFrame
	Mods   *ReplayMods
}

func (replay *ReplayData) String() string {
	return fmt.Sprintf(
		"ReplayData{%d frames}", // TODO: Header
		len(replay.Frames),
	)
}

func (replay *ReplayData) Serialize(stream *IOStream) {
	replay.SerializeHeader(stream)
	replay.SerializeFrames(stream)
	replay.Mods.Serialize(stream)
}

func (replay *ReplayData) SerializeHeader(stream *IOStream) {
	stream.Write([]byte{0xf0, 0xa0, 0xc0, 0xe0})
	stream.WriteI8(replay.Mode)
	stream.WriteI32(int32(replay.ReplayVersion))
	stream.WriteString(replay.BeatmapChecksum)
	stream.WriteString(replay.PlayerName)
	stream.WriteString(replay.ScoreChecksum)
	stream.WriteU32(replay.Count300)
	stream.WriteU32(replay.Count100)
	stream.WriteU32(replay.Count50)
	stream.WriteU32(replay.CountGeki)
	stream.WriteU32(replay.CountGood)
	stream.WriteU32(replay.CountMiss)
	stream.WriteF64(replay.TotalScore)
	stream.WriteU32(replay.MaxCombo)
	stream.WriteBool(replay.FullCombo)
	stream.WriteReplayDateTime(replay.Time, replay.TimeSpec)
}

func (replay *ReplayData) SerializeFrames(stream *IOStream) {
	replayStream := NewIOStream([]byte{}, binary.BigEndian)
	replayStream.WriteU32(uint32(len(replay.Frames)))

	for _, frame := range replay.Frames {
		frame.Serialize(replayStream)
	}

	decompressed := replayStream.Get()
	compressed := bytes.NewBuffer([]byte{})

	zlibWriter := zlib.NewWriter(compressed)
	zlibWriter.Write(decompressed)
	zlibWriter.Close()

	frameData := NewIOStream([]byte{}, binary.BigEndian)
	frameData.WriteU32(uint32(len(decompressed)))
	frameData.Write(compressed.Bytes())

	stream.WriteQByteArray(frameData.Get())
}

func (replay *ReplayData) Accuracy() float64 {
	totalHits := replay.Count300 + replay.Count100 + replay.Count50
	return float64(replay.Count300*300+replay.Count100*100+replay.Count50*50) / float64(totalHits*300)
}

func (replay *ReplayData) Grade() Grade {
	return CalculateGrade(
		true,
		int(replay.Count300),
		int(replay.Count100),
		int(replay.Count50),
		int(replay.CountMiss),
		replay.Mods.Hidden,
	)
}

func (replay *ReplayData) Checksum() string {
	return CreateReplayChecksum(
		replay.PlayerName,
		replay.BeatmapChecksum,
		true,
		int(replay.Count300),
		int(replay.Count100),
		int(replay.Count50),
		int(replay.CountGeki),
		int(replay.CountGood),
		int(replay.CountMiss),
		int(replay.MaxCombo),
		replay.FullCombo,
		int(math.Round(replay.TotalScore)),
		replay.Grade(),
		CreateModsChecksumToken(
			replay.Mods.ArOffset,
			replay.Mods.OdOffset,
			replay.Mods.CsOffset,
			replay.Mods.HpOffset,
			replay.Mods.PsOffset,
			replay.Mods.Hidden,
			replay.Mods.NoFail,
			replay.Mods.Autoplay,
		),
	)
}

type ReplayMods struct {
	ArOffset int
	OdOffset int
	CsOffset int
	HpOffset int
	PsOffset int
	Hidden   bool
	NoFail   bool
	Autoplay bool
}

func (mods *ReplayMods) Serialize(stream *IOStream) {
	if mods == nil {
		mods = &ReplayMods{}
	}

	stream.WriteI32(int32(mods.ArOffset))
	stream.WriteI32(int32(mods.OdOffset))
	stream.WriteI32(int32(mods.CsOffset))
	stream.WriteI32(int32(mods.HpOffset))
	stream.WriteI32(int32(mods.PsOffset))
	stream.WriteBool(mods.Hidden)
	stream.WriteBool(mods.NoFail)
	stream.WriteBool(mods.Autoplay)
}

func ReadReplayMods(stream *IOStream) *ReplayMods {
	return &ReplayMods{
		ArOffset: int(stream.ReadI32()),
		OdOffset: int(stream.ReadI32()),
		CsOffset: int(stream.ReadI32()),
		HpOffset: int(stream.ReadI32()),
		PsOffset: int(stream.ReadI32()),
		Hidden:   stream.ReadBool(),
		NoFail:   stream.ReadBool(),
		Autoplay: stream.ReadBool(),
	}
}

type ReplayFrame struct {
	Time        uint32
	MouseX      float64
	MouseY      float64
	ButtonState uint32
}

func (frame *ReplayFrame) String() string {
	return FormatStruct(frame)
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

func ReadReplayFrames(stream *IOStream) (frames []*ReplayFrame, err error) {
	defer HandlePanic(&err)

	if stream.Available() < 4 {
		return frames, fmt.Errorf("replay is too short")
	}

	expectedSize := stream.ReadU32()
	compressedReplayData := stream.ReadAll()

	reader := bytes.NewReader(compressedReplayData)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return frames, err
	}

	replayData, err := io.ReadAll(zlibReader)
	if err != nil {
		return frames, err
	}

	if expectedSize != 0 && uint32(len(replayData)) != expectedSize {
		return frames, fmt.Errorf(
			"unexpected replay size: got %d bytes, expected %d",
			len(replayData),
			expectedSize,
		)
	}

	if len(replayData) < 4 {
		return frames, fmt.Errorf("replay data is too short")
	}

	stream = NewIOStream(replayData, binary.BigEndian)
	frameAmount := stream.ReadU32()
	frames = make([]*ReplayFrame, frameAmount)

	// One frame is 24 bytes
	expectedFrameBytes := 24 * frameAmount

	// Check if we have enough data for all frames
	if stream.Available() < int(expectedFrameBytes) {
		return frames, fmt.Errorf(
			"not enough data for %d frames, got %d bytes",
			frameAmount, stream.Available(),
		)
	}

	for i := range frameAmount {
		frames[i] = ReadReplayFrame(stream)
	}
	return frames, nil
}

func ReadFullReplay(stream *IOStream) (*ReplayData, error) {
	magicBytes := stream.Read(4)

	if !bytes.Equal(magicBytes, []byte{0xf0, 0xa0, 0xc0, 0xe0}) {
		return nil, fmt.Errorf("invalid replay format: missing magic bytes")
	}

	replayData := &ReplayData{}
	replayData.Mode = stream.ReadI8()
	replayData.ReplayVersion = int(stream.ReadI32())
	replayData.BeatmapChecksum = stream.ReadString()
	replayData.PlayerName = stream.ReadString()
	replayData.ScoreChecksum = stream.ReadString()
	replayData.Count300 = stream.ReadU32()
	replayData.Count100 = stream.ReadU32()
	replayData.Count50 = stream.ReadU32()
	replayData.CountGeki = stream.ReadU32()
	replayData.CountGood = stream.ReadU32()
	replayData.CountMiss = stream.ReadU32()
	replayData.TotalScore = stream.ReadF64()
	replayData.MaxCombo = stream.ReadU32()
	replayData.FullCombo = stream.ReadBool()
	replayData.Time, replayData.TimeSpec = stream.ReadReplayDateTime()

	replayBytes, err := stream.ReadQByteArray()
	if err != nil {
		return replayData, err
	}

	replayData.Mods = ReadReplayMods(stream)
	replayData.Frames, err = ReadReplayFrames(NewIOStream(replayBytes, binary.BigEndian))
	return replayData, err
}
