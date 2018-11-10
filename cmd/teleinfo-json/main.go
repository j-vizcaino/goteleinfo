package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/j-vizcaino/goteleinfo"
)

func main() {
	var serialDevice string

	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.Parse()

	port, err := teleinfo.OpenPort(serialDevice)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer port.Close()
	reader := teleinfo.NewReader(port)

	for {
		frame, err := reader.ReadFrame()
		if err != nil {
			fmt.Printf("Error reading frame from '%s' (%s)\n", serialDevice, err)
			continue
		}
		doc, _ := json.Marshal(frame)
		fmt.Println(string(doc))
	}
}
