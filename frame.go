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
	// Mode returns the mode of frame (historic or standard)
	Mode() string
	// GetStringField returns the value of a field as a string.
	GetStringField(string) (string, bool)
	// GetUIntField returns type value of a field as an unsigned integer.
	GetUIntField(string) (uint, bool)
	// GetMap returns the raw map object
	GetMap() map[string]string
}

type frame map[string]string

func (f frame) Type() string {
	return f["OPTARIF"]
}

func (f frame) Mode() string {
	if _, ok := f["OPTARIF"]; ok {
		return "historic"
	}
	return "standard"
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

func (f frame) GetMap() map[string]string {
	return f
}

// NOTES:
// * used by decodeFrame()
// * []byte cannot be const :(
var (
	historicFieldSeparator = []byte("\r\n")
	historicEltSeparator   = []byte(" ")
	standardFieldSeparator = []byte("\r")
	standardEltSeparator   = []byte("\t")
)

func decodeStandardFrame(rawFrame []byte) (Frame, error) {
	const (
		checksumLength = 1
	)
	var (
		name      []byte
		timestamp []byte
		value     []byte
		trail     []byte
	)
	fields := bytes.Split(rawFrame, standardFieldSeparator)
	info := frame{}
	for _, field := range fields {
		field = bytes.TrimPrefix(field, []byte("\n"))
		elts := bytes.Split(field, standardEltSeparator)

		switch len(elts) {
		case 1:
			break
		case 3:
			timestamp = []byte("")
			name, value, trail = elts[0], elts[1], elts[2]
		case 4:
			name, timestamp, value, trail = elts[0], elts[1], elts[2], elts[3]
		default:
			incrementErrorCounter(frameDecodeErrorCounter, "invalid_field")
			return nil, fmt.Errorf("error decoding frame, invalid number of elements for data (data: '%s' len:%d)", field, len(elts))
		}

		if len(trail) != checksumLength {
			incrementErrorCounter(frameDecodeErrorCounter, "invalid_checksum_length")
			return nil, fmt.Errorf("error decoding frame, invalid checksum length (actual: %d, expected: %d)", len(trail), checksumLength)
		}
		readChecksum := byte(trail[0])
		expectedChecksum := standardChecksum(name, timestamp, value)
		if readChecksum != expectedChecksum {
			incrementErrorCounter(frameDecodeErrorCounter, "checksum_error")
			return nil, fmt.Errorf("error decoding frame, invalid checksum (field: '%s', value: '%s' read: '%c', expected: '%c')", name, value, readChecksum, expectedChecksum)
		}
		info[string(name)] = string(value)
	}
	frameDecodedCounter.Inc()
	return info, nil
}

func decodeHistoricFrame(rawFrame []byte) (Frame, error) {
	const (
		checksumLength = 1
	)

	strFrame := bytes.Trim(rawFrame, "\r\n")

	fields := bytes.Split(strFrame, historicFieldSeparator)
	info := frame{}
	for _, field := range fields {
		elts := bytes.SplitN(field, historicEltSeparator, 3)

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
		expectedChecksum := historicChecksum(name, value)
		if readChecksum != expectedChecksum {
			incrementErrorCounter(frameDecodeErrorCounter, "checksum_error")
			return nil, fmt.Errorf("error decoding frame, invalid checksum (field: '%s', value: '%s', read: '%c', expected: '%c')", name, value, readChecksum, expectedChecksum)
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

func standardChecksum(name []byte, timestamp []byte, value []byte) byte {
	// NOTE: 0x09 == ASCII horizontal tab
	// https://www.enedis.fr/sites/default/files/Enedis-NOI-CPT_54E.pdf - Version 3 - 01/06/2018 - Page 14/38
	checksum := sum(name) + byte(0x09) + sum(value) + byte(0x09)
	var b []byte
	if !bytes.Equal(timestamp, b) {
		checksum = sum(name) + byte(0x09) + sum(timestamp) + byte(0x09) + sum(value) + byte(0x09)
	}

	// S1 AND 0x3F + 0x20
	checksum = (checksum & 0x3F) + 0x20
	return checksum
}

func historicChecksum(name []byte, value []byte) byte {
	// NOTE: 0x20 == ASCII space char
	checksum := sum(name) + byte(0x20) + sum(value)

	// Map to a single char E [0x20;0x7F]
	checksum = (checksum & 0x3F) + 0x20
	return checksum
}
