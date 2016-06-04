package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"bufio"
	"github.com/tarm/serial"
	"github.com/golang/glog"
	"strconv"
	"time"
	"bytes"
)

func ReadFrame(r *bufio.Reader) (frame []byte, err error) {
	const (
		FRAME_START byte = 0x2
		FRAME_END byte = 0x3
	)
	// TODO: check for overflow
	// TODO: check for interrupted frame marker
	if _, err = r.ReadSlice(FRAME_START); err != nil {
		return
	}
	if frame, err = r.ReadBytes(FRAME_END); err != nil {
		return
	}
	frame = frame[0:len(frame)-1]
	return
}

type ConsoInfo struct {
	Timestamp time.Time
	PAPP_va   uint32
	HC_wh     uint32
	HP_wh     uint32
	IsHP      bool
}

func ExtractNumber(value string) uint32 {
	num, _ := strconv.ParseUint(value, 10, 32)
	return uint32(num)
}

func sum(a []byte) (res byte) {
	res = 0
	for _, c := range a {
		res += c
	}
	return
}

func ComputeChecksum(name []byte, value []byte) byte {
	// NOTE: 0x20 == ASCII space char
	checksum := sum(name) + byte(0x20) + sum(value)

	// Map to a single char E [0x20;0x7F]
	checksum = (checksum & 0x3F) + 0x20
	return checksum
}

func ExtractInfo(fields map[string]string) ConsoInfo {

	return ConsoInfo{
		Timestamp: time.Now(),
		PAPP_va:   ExtractNumber(fields["PAPP"]),
		HC_wh:     ExtractNumber(fields["HCHC"]),
		HP_wh:     ExtractNumber(fields["HCHP"]),
		IsHP:      fields["PTEC"] == "HP..",
	}
}

func DecodeFrame(rawFrame []byte) (map[string]string, error) {
	const (
		CHECKSUM_LENGTH = 1
	)
	FIELD_SEPARATOR := []byte("\r\n")
	ELTS_SEPARATOR := []byte(" ")

	strFrame := bytes.Trim(rawFrame, "\r\n")

	fields := bytes.Split(strFrame, FIELD_SEPARATOR)
	info := make(map[string]string)
	for _, field := range fields {
		elts := bytes.SplitN(field, ELTS_SEPARATOR, 3)

		if len(elts) != 3 {
			return nil, fmt.Errorf("invalid number of elements for data (data: '%s')", field)
		}
		name, value, trail := elts[0], elts[1], elts[2]

		if len(trail) != CHECKSUM_LENGTH {
			return nil, fmt.Errorf("invalid checksum length (actual: %d, expected: %d)", len(trail), CHECKSUM_LENGTH)
		}
		readChecksum := byte(trail[0])
		computedChecksum := ComputeChecksum(name, value)
		if readChecksum != computedChecksum {
			return nil, fmt.Errorf("invalid checksum (field: '%s', value: '%s', read: '%c', expected: '%c'", name, value, readChecksum, computedChecksum)
		}
		info[string(name)] = string(value)
	}
	return info, nil
}

func main() {
	var serialDevice string

	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.Parse()

	cfg := &serial.Config{
		Name:     serialDevice,
		Baud:     1200,
		Size:     7,
		Parity:   serial.ParityEven,
		StopBits: serial.Stop1,
	}

	port, err := serial.OpenPort(cfg)
	if err != nil {
		glog.Exitf("Error opening device '%s' (%s)\n", serialDevice, err)
	}
	defer port.Close()

	reader := bufio.NewReader(port)
	for {
		frame, err := ReadFrame(reader)
		if err != nil {
			glog.Exitf("Error reading frame from '%s' (%s)\n", serialDevice, err)
		}
		fields, err := DecodeFrame(frame)
		if err != nil {
			glog.Warningf("Error decoding frame: %s. Frame has been dropped.", err)
			continue
		}
		info := ExtractInfo(fields)
		doc, err := json.Marshal(info)
		fmt.Println(string(doc))
	}
}
