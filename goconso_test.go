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
	INVALID_CHECKSUM := []byte("\rHCHC 09876543 e\n")
	INVALID_CHECKSUM_SIZE := []byte("\rHCHC 09876543 ea\n")
	RAW_FRAME := []byte("\rPAPP 012345  \r\nPTEC HP..  \r\nHCHC 09876543 @\r\nHCHP 1654800 K\n")
	EXPECTED_MAP := map[string]string{
		"PAPP": "012345",
		"PTEC": "HP..",
		"HCHC": "09876543",
		"HCHP": "1654800",
	}

	r, err := DecodeFrame(INVALID_CHECKSUM)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = DecodeFrame(INVALID_CHECKSUM_SIZE)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = DecodeFrame(RAW_FRAME)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_MAP, r)
}

func TestComputeChecksum(t *testing.T) {
	NAME := []byte("PAPP")
	VALUE := []byte("012345")
	EXPECTED_CHECKSUM := byte(0x20)

	assert.Equal(t, EXPECTED_CHECKSUM, ComputeChecksum(NAME, VALUE))
}