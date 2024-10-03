package common

import "encoding/binary"

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

func (stream *IOStream) ReadString() string {
	length := stream.ReadU32()

	if length == 0 {
		return ""
	}

	data := stream.Read(int(length))
	chars := make([]rune, 0, length)

	for _, char := range data {
		if char == 0 {
			continue
		}

		chars = append(chars, rune(char))
	}

	return string(chars)
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

func (stream *IOStream) WriteString(value string) {
	stream.WriteU32(uint32(len(value) * 2))

	for _, c := range value {
		stream.Write([]byte{
			0x00,
			byte(c),
		})
	}

	stream.Write([]byte{0x00})
}
