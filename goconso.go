package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"goconso/teleinfo"

	"github.com/golang/glog"
)

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

func main() {
	var serialDevice string

	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.Parse()

	port, err := teleinfo.OpenPort(serialDevice)
	if err != nil {
		glog.Exitf("Error: %s", err)
	}
	defer port.Close()
	reader := teleinfo.NewReader(port)

	for {
		frame, err := reader.ReadFrame()
		if err != nil {
			glog.Exitf("Error reading frame from '%s' (%s)\n", serialDevice, err)
		}
		info := ExtractInfo(frame)
		doc, err := json.Marshal(info)
		fmt.Println(string(doc))
	}
}
