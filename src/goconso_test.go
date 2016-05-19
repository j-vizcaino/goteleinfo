package main
import (
	"testing"
	"bytes"
	"bufio"
	"github.com/stretchr/testify/assert"
)

func TestReadFrame(t *testing.T) {
	EXPECTED_FRAME := []byte("frame content")
	BUFFER := bytes.NewBufferString("some junk prefix\r\n")
	BUFFER.WriteByte(0x2)
	BUFFER.Write(EXPECTED_FRAME)
	BUFFER.WriteByte(0x3)
	BUFFER.WriteString("trailing junk")

	buf := bytes.NewReader(BUFFER.Bytes())
	reader := bufio.NewReader(buf)

	frame, err := ReadFrame(reader)
	assert.Nil(t, err)
	assert.Equal(t, frame, EXPECTED_FRAME)
}

func TestDecodeFrame(t *testing.T) {
	RAW_FRAME := []byte("\rPAPP 012345 x\r\nPTEC HP.. u\r\nHCHC 09876543 e\r\nHCHP 1654800 u\n")
	EXPECTED_MAP := map[string]string{
		"PAPP": "012345",
		"PTEC": "HP..",
		"HCHC": "09876543",
		"HCHP": "1654800",
	}

	r := DecodeFrame(RAW_FRAME)
	assert.Equal(t, r, EXPECTED_MAP)
}