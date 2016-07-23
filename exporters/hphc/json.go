package hphc

import (
	"encoding/json"
	"fmt"
	"goconso/exporters"
	"goconso/teleinfo"
)

func init() {
	exporters.Register("hphc.json", NewJSONExporter)
}

type JSONExporter struct{}

func NewJSONExporter() teleinfo.Exporter {
	return &JSONExporter{}
}

func (x *JSONExporter) ExportFrame(f teleinfo.Frame) error {
	record := NewRecord(f)
	doc, err := json.Marshal(record)
	if err != nil {
		return err
	}
	fmt.Println(string(doc))
	return nil
}
