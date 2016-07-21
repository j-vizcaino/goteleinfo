package teleinfo

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRawFrame(t *testing.T) {
	expectedFrame := []byte("frame content")
	bufferContent := bytes.NewBufferString("some junk prefix\r\n")
	bufferContent.WriteByte(0x2)
	bufferContent.Write(expectedFrame)
	bufferContent.WriteByte(0x3)
	bufferContent.WriteString("trailing junk")

	buf := bytes.NewReader(bufferContent.Bytes())
	reader := bufio.NewReader(buf)

	frame, err := readRawFrame(reader)
	assert.Nil(t, err)
	assert.Equal(t, frame, expectedFrame)
}

var RawFramePayload = []byte("\rPAPP 012345  \r\nPTEC HP..  \r\nHCHC 09876543 @\r\nHCHP 1654800 K\n")
var RawFrame = []byte("\x02\rPAPP 012345  \r\nPTEC HP..  \r\nHCHC 09876543 @\r\nHCHP 1654800 K\n\x03")
var DecodedFrame = Frame{
	"PAPP": "012345",
	"PTEC": "HP..",
	"HCHC": "09876543",
	"HCHP": "1654800",
}

func TestReadFrame(t *testing.T) {
	reader := NewReader(bytes.NewReader(RawFrame))

	frame, err := reader.ReadFrame()
	assert.Nil(t, err)
	assert.Equal(t, DecodedFrame, frame)
}

func TestDecodeFrame(t *testing.T) {
	InvalidChecksumValue := []byte("\rHCHC 09876543 e\n")
	InvalidChecksumLength := []byte("\rHCHC 09876543 ea\n")
	InvalidFieldCount := []byte("\rHCHC 09876543\n")

	r, err := decodeFrame(InvalidChecksumValue)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(InvalidChecksumLength)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(InvalidFieldCount)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(RawFramePayload)
	assert.Nil(t, err)
	assert.Equal(t, DecodedFrame, r)
}

func TestComputeChecksum(t *testing.T) {
	name := []byte("PAPP")
	value := []byte("012345")
	expectedChecksum := byte(0x20)

	assert.Equal(t, expectedChecksum, computeChecksum(name, value))
}
