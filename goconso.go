package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"bufio"
	"github.com/tarm/serial"
	"strconv"
	"time"
	"strings"
)

const FIELD_SEPARATOR = "\r\n"
const ELTS_SEPARATOR = " "

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

func ExtractInfo(fields map[string]string) ConsoInfo {

	return ConsoInfo{
		Timestamp: time.Now(),
		PAPP_va:   ExtractNumber(fields["PAPP"]),
		HC_wh:     ExtractNumber(fields["HCHC"]),
		HP_wh:     ExtractNumber(fields["HCHP"]),
		IsHP:      fields["PTEC"] == "HP..",
	}
}

func DecodeFrame(rawFrame []byte) map[string]string {
	strFrame := strings.Trim(string(rawFrame), "\r\n")

	fields := strings.Split(strFrame, FIELD_SEPARATOR)
	info := make(map[string]string)
	for _, field := range fields {
		elts := strings.SplitN(field, ELTS_SEPARATOR, 3)
		// TODO: handle incorrect number of fields
		name, value, _ := elts[0], elts[1], elts[2]
		// TODO: verify checksum
		info[name] = value
	}
	return info
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
		fmt.Printf("Error opening device '%s' (%s)\n", serialDevice, err)
		return
	}
	defer port.Close()

	reader := bufio.NewReader(port)
	for {
		frame, err := ReadFrame(reader)
		if err != nil {
			fmt.Printf("Error reading frame from '%s' (%s)\n", serialDevice, err)
			return
		}
		fields := DecodeFrame(frame)
		info := ExtractInfo(fields)
		doc, err := json.Marshal(info)
		fmt.Println(string(doc))
	}
}
