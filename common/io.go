package common

import (
	"encoding/binary"
	"math"
	"time"
)

type IOStream struct {
	data     []byte
	position int
	endian   binary.ByteOrder
}

func NewIOStream(data []byte, endian binary.ByteOrder) *IOStream {
	return &IOStream{
		data:     data,
		position: 0,
		endian:   endian,
	}
}

func (stream *IOStream) Push(data []byte) {
	stream.data = append(stream.data, data...)
}

func (stream *IOStream) Get() []byte {
	return stream.data
}

func (stream *IOStream) Len() int {
	return len(stream.data)
}

func (stream *IOStream) Available() int {
	return stream.Len() - stream.position
}

func (stream *IOStream) Tell() int {
	return stream.position
}

func (stream *IOStream) Seek(position int) {
	stream.position = position
}

func (stream *IOStream) Skip(offset int) {
	stream.position += offset
}

func (stream *IOStream) Eof() bool {
	return stream.position >= stream.Len()
}

func (stream *IOStream) Read(size int) []byte {
	if stream.Eof() {
		return []byte{}
	}

	if stream.Available() < size {
		size = stream.Available()
	}

	data := stream.data[stream.position : stream.position+size]
	stream.position += size

	return data
}

func (stream *IOStream) ReadAll() []byte {
	return stream.Read(stream.Available())
}

func (stream *IOStream) ReadU8() uint8 {
	return stream.Read(1)[0]
}

func (stream *IOStream) ReadU16() uint16 {
	return stream.endian.Uint16(stream.Read(2))
}

func (stream *IOStream) ReadU32() uint32 {
	return stream.endian.Uint32(stream.Read(4))
}

func (stream *IOStream) ReadU64() uint64 {
	return stream.endian.Uint64(stream.Read(8))
}

func (stream *IOStream) ReadI8() int8 {
	return int8(stream.Read(1)[0])
}

func (stream *IOStream) ReadI16() int16 {
	return int16(stream.ReadU16())
}

func (stream *IOStream) ReadI32() int32 {
	return int32(stream.ReadU32())
}

func (stream *IOStream) ReadI64() int64 {
	return int64(stream.ReadU64())
}

func (stream *IOStream) ReadF32() float32 {
	bits := stream.ReadU32()
	return math.Float32frombits(bits)
}

func (stream *IOStream) ReadF64() float64 {
	bits := stream.ReadU64()
	return math.Float64frombits(bits)
}

func (stream *IOStream) ReadBool() bool {
	return stream.ReadU8() == 1
}

func (stream *IOStream) ReadString() string {
	length := stream.ReadU32()

	if length == 0 {
		return ""
	}

	data := stream.Read(int(length))
	chars := make([]rune, 0, length)

	for i := 0; i < len(data); i += 2 {
		char := rune(stream.endian.Uint16(data[i : i+2]))
		chars = append(chars, char)
	}

	return string(chars)
}

func (stream *IOStream) ReadIntList() []uint32 {
	length := stream.ReadU32()

	if length == 0 {
		return []uint32{}
	}

	list := make([]uint32, 0, length)

	for range length {
		list = append(list, stream.ReadU32())
	}

	return list
}

func (stream *IOStream) ReadDateTime() time.Time {
	// Convert julian date to time.Time
	jd := float64(stream.ReadI32())
	time := JulianToTime(jd)
	return time
}

func (stream *IOStream) Write(data []byte) {
	stream.Push(data)
}

func (stream *IOStream) WriteU8(value uint8) {
	stream.Write([]byte{value})
}

func (stream *IOStream) WriteU16(value uint16) {
	data := make([]byte, 2)
	stream.endian.PutUint16(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteU32(value uint32) {
	data := make([]byte, 4)
	stream.endian.PutUint32(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteU64(value uint64) {
	data := make([]byte, 8)
	stream.endian.PutUint64(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteI8(value int8) {
	stream.WriteU8(uint8(value))
}

func (stream *IOStream) WriteI16(value int16) {
	stream.WriteU16(uint16(value))
}

func (stream *IOStream) WriteI32(value int32) {
	stream.WriteU32(uint32(value))
}

func (stream *IOStream) WriteI64(value int64) {
	stream.WriteU64(uint64(value))
}

func (stream *IOStream) WriteF32(value float32) {
	bits := math.Float32bits(value)
	stream.WriteU32(bits)
}

func (stream *IOStream) WriteF64(value float64) {
	bits := math.Float64bits(value)
	stream.WriteU64(bits)
}

func (stream *IOStream) WriteBool(value bool) {
	if value {
		stream.WriteU8(1)
	} else {
		stream.WriteU8(0)
	}
}

func (stream *IOStream) WriteString(value string) {
	stream.WriteU32(uint32(len(value) * 2))

	for _, c := range value {
		stream.WriteU16(uint16(c))
	}
}

func (stream *IOStream) WriteIntList(list []uint32) {
	stream.WriteU32(uint32(len(list)))

	for _, value := range list {
		stream.WriteU32(value)
	}
}

func (stream *IOStream) WriteDateTime(value time.Time) {
	// Convert time.Time to julian date
	jd := TimeToJulian(value)
	stream.WriteI32(int32(jd))
}
