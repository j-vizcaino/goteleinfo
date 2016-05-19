package teleinfo

import (
	"bufio"
	"fmt"
	"io"
)

type Reader interface {
	ReadFrame() (Frame, error)
}

type reader struct {
	buffer *bufio.Reader
}

func NewReader(r io.Reader) Reader {
	return &reader{
		buffer: bufio.NewReader(r),
	}
}


func readRawFrame(buffer *bufio.Reader) ([]byte, error) {
	const (
		FrameStart byte = 0x2
		FrameEnd   byte = 0x3
	)
	var frame []byte
	var err error
	// TODO: check for overflow
	// TODO: check for interrupted frame marker
	if _, err = buffer.ReadSlice(FrameStart); err != nil {
		return nil, fmt.Errorf("error looking for start of frame marker (%s)", err)
	}
	if frame, err = buffer.ReadBytes(FrameEnd); err != nil {
		return nil, fmt.Errorf("error looking for end of frame marker (%s)", err)
	}

	if len(frame) == 0 {
		return frame, fmt.Errorf("read empty frame")
	}
	frame = frame[0 : len(frame)-1]
	return frame, nil
}

func (t *reader) ReadFrame() (Frame, error) {
	rawFrame, err := readRawFrame(t.buffer)
	if err != nil {
		return nil, err
	}
	return decodeFrame(rawFrame)
}
