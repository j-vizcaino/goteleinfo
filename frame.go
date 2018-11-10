package teleinfo

import (
	"bytes"
	"fmt"
	"strconv"
)

// Frame holds a single Teleinfo frame.
type Frame interface {
	// Type returns the type of frame (see `OPTARIF` field)
	Type() string
	// GetStringField returns the value of a field as a string.
	GetStringField(string) (string, bool)
	// GetUIntField returns type value of a field as an unsigned integer.
	GetUIntField(string) (uint, bool)
}

type frame map[string]string

func (f frame) Type() string {
	return f["OPTARIF"]
}

func (f frame) GetStringField(name string) (string, bool) {
	v, ok := f[name]
	return v, ok
}

func (f frame) GetUIntField(name string) (uint, bool) {
	s, ok := f[name]
	if !ok {
		return 0, ok
	}
	num, err := strconv.ParseUint(s, 10, 32)
	ok = err == nil
	return uint(num), ok
}

// NOTES:
// * used by decodeFrame()
// * []byte cannot be const :(
var fieldSeparator = []byte("\r\n")
var eltSeparator = []byte(" ")

func decodeFrame(rawFrame []byte) (Frame, error) {
	const (
		checksumLength = 1
	)

	strFrame := bytes.Trim(rawFrame, "\r\n")

	fields := bytes.Split(strFrame, fieldSeparator)
	info := frame{}
	for _, field := range fields {
		elts := bytes.SplitN(field, eltSeparator, 3)

		if len(elts) != 3 {
			incrementErrorCounter(frameDecodeErrorCounter, "invalid_field")
			return nil, fmt.Errorf("error decoding frame, invalid number of elements for data (data: '%s')", field)
		}
		name, value, trail := elts[0], elts[1], elts[2]

		if len(trail) != checksumLength {
			incrementErrorCounter(frameDecodeErrorCounter, "invalid_checksum_length")
			return nil, fmt.Errorf("error decoding frame, invalid checksum length (actual: %d, expected: %d)", len(trail), checksumLength)
		}
		readChecksum := byte(trail[0])
		expectedChecksum := computeChecksum(name, value)
		if readChecksum != expectedChecksum {
			incrementErrorCounter(frameDecodeErrorCounter, "checksum_error")
			return nil, fmt.Errorf("error decoding frame, invalid checksum (field: '%s', value: '%s', read: '%c', expected: '%c'", name, value, readChecksum, expectedChecksum)
		}
		info[string(name)] = string(value)
	}
	frameDecodedCounter.Inc()
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
