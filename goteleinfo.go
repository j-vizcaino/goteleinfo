package main

import (
	"flag"
	"fmt"
	"goteleinfo/exporters"
	_ "goteleinfo/exporters/hphc"
	"goteleinfo/teleinfo"

	"github.com/golang/glog"
)

func main() {
	var serialDevice string
	var exporterName string

	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.StringVar(&exporterName, "export", "", fmt.Sprintf("Exporter module name. Available: %v", exporters.List()))
	flag.Parse()

	newExporter := exporters.FindByName(exporterName)
	if newExporter == nil {
		glog.Exitf("Unknown exporter '%s', valid choices are %v\n", exporterName, exporters.List())
	}
	exporter, err := newExporter()
	if err != nil {
		glog.Exitf("Error creating exporter '%s' (%s)\n", exporterName, err)
	}
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
		err = exporter.ExportFrame(frame)
		if err != nil {
			glog.Errorf("Error exporting frame (%s)\n", err)
		}
	}
}
