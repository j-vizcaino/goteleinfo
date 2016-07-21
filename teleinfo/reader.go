package teleinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Frame map[string]string

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

// NOTES:
// * used be decodeFrame()
// * []byte cannot be const :(
var fieldSeparator = []byte("\r\n")
var eltSeparator = []byte(" ")

func decodeFrame(rawFrame []byte) (Frame, error) {
	const (
		checksumLength = 1
	)

	strFrame := bytes.Trim(rawFrame, "\r\n")

	fields := bytes.Split(strFrame, fieldSeparator)
	info := make(map[string]string)
	for _, field := range fields {
		elts := bytes.SplitN(field, eltSeparator, 3)

		if len(elts) != 3 {
			return nil, fmt.Errorf("error decoding frame, invalid number of elements for data (data: '%s')", field)
		}
		name, value, trail := elts[0], elts[1], elts[2]

		if len(trail) != checksumLength {
			return nil, fmt.Errorf("error decoding frame, invalid checksum length (actual: %d, expected: %d)", len(trail), checksumLength)
		}
		readChecksum := byte(trail[0])
		expectedChecksum := computeChecksum(name, value)
		if readChecksum != expectedChecksum {
			return nil, fmt.Errorf("error decoding frame, invalid checksum (field: '%s', value: '%s', read: '%c', expected: '%c'", name, value, readChecksum, expectedChecksum)
		}
		info[string(name)] = string(value)
	}
	return info, nil
}

func sum(a []byte) (res byte) {
	res = 0
	for _, c := range a {
		res += c
	}
	return
}

func computeChecksum(name []byte, value []byte) byte {
	// NOTE: 0x20 == ASCII space char
	checksum := sum(name) + byte(0x20) + sum(value)

	// Map to a single char E [0x20;0x7F]
	checksum = (checksum & 0x3F) + 0x20
	return checksum
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
		return Frame{}, err
	}
	return decodeFrame(rawFrame)
}
