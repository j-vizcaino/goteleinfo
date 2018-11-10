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
var DecodedFrame = frame{
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

func TestReadFrameNoStart(t *testing.T) {
	invalidData := []byte("qwertyuiop")
	reader := NewReader(bytes.NewReader(invalidData))

	_, err := reader.ReadFrame()
	assert.Error(t, err)
}

func TestReadFrameNoEnd(t *testing.T) {
	invalidData := make([]byte, 0, bufio.MaxScanTokenSize+1)
	// Set frame start
	invalidData = append(invalidData, byte(0x2))
	for i := 1; i < cap(invalidData); i++ {
		invalidData = append(invalidData, byte(i%255))
	}

	reader := NewReader(bytes.NewReader(invalidData))
	_, err := reader.ReadFrame()
	assert.Error(t, err)
}
