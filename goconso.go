package main

import (
	"flag"

	"goconso/exporters/hphc"
	"goconso/teleinfo"

	"github.com/golang/glog"
)

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
	exporter := hphc.NewJSONExporter()

	for {
		frame, err := reader.ReadFrame()
		if err != nil {
			glog.Exitf("Error reading frame from '%s' (%s)\n", serialDevice, err)
		}
		exporter.ExportFrame(frame)
	}
}
