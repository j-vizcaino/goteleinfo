package teleinfo

import (
	"bufio"
	"fmt"
	"io"
)

// Reader defines an interface to read raw Teleinfo frames.
type Reader interface {
	// ReadFrame reads a raw Teleinfo frame.
	ReadFrame() (Frame, error)
}

type reader struct {
	buffer *bufio.Reader
}

// NewReader create a Teleinfo frame reader from a simple Reader.
// r is usually the struct returned by OpenPort().
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
		incrementErrorCounter(frameReadErrorCounter, "no_frame_start_marker")
		return nil, fmt.Errorf("error looking for start of frame marker (%s)", err)
	}
	if frame, err = buffer.ReadBytes(FrameEnd); err != nil {
		incrementErrorCounter(frameReadErrorCounter, "no_frame_end_marker")
		return nil, fmt.Errorf("error looking for end of frame marker (%s)", err)
	}

	if len(frame) == 0 {
		incrementErrorCounter(frameReadErrorCounter, "empty_frame")
		return frame, fmt.Errorf("read empty frame")
	}
	frame = frame[0 : len(frame)-1]
	frameReadCounter.Inc()
	return frame, nil
}

func (t *reader) ReadFrame() (Frame, error) {
	rawFrame, err := readRawFrame(t.buffer)
	if err != nil {
		return nil, err
	}
	return decodeFrame(rawFrame)
}
