package hphc

import (
	"encoding/json"
	"fmt"
	"goteleinfo/exporters"
	"goteleinfo/teleinfo"
)

func init() {
	exporters.Register("hphc.json", newJSONExporter)
}

type jsonExporter struct{}

func newJSONExporter() (teleinfo.Exporter, error) {
	return &jsonExporter{}, nil
}

func (x *jsonExporter) ExportFrame(f teleinfo.Frame) error {
	rec := newRecord(f)
	doc, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	fmt.Println(string(doc))
	return nil
}
