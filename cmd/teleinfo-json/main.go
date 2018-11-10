package main

import (
	"container/ring"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/j-vizcaino/goteleinfo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"sync"
)

func readFrames(reader teleinfo.Reader, framesChan chan<- teleinfo.Frame) {
	for {
		frame, err := reader.ReadFrame()
		if err != nil {
			fmt.Printf("Error reading Teleinfo frame: %s\n", err)
			continue
		}
		framesChan <- frame
	}
}

func main() {
	var serialDevice string
	var listenAddress string
	var framesCount int

	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.StringVar(&listenAddress, "listen-address", "localhost:9000", "HTTP service listen address")
	flag.IntVar(&framesCount, "frames-count", 20, "Number of Teleinfo frames to serve under /frames")
	flag.Parse()

	port, err := teleinfo.OpenPort(serialDevice)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer port.Close()

	framesChan := make(chan teleinfo.Frame, 10)
	framesBuffer := ring.New(framesCount)
	mutex := &sync.Mutex{}

	// Read Teleinfo frames and send them into framesChan
	go readFrames(teleinfo.NewReader(port), framesChan)

	// Enqueue teleinfo.Frame into a fixed-length ring buffer
	go func() {
		for frame := range framesChan {
			mutex.Lock()
			framesBuffer.Value = frame
			framesBuffer = framesBuffer.Next()
			mutex.Unlock()
		}
	}()

	http.HandleFunc("/frames", func(w http.ResponseWriter, req *http.Request) {
		// Convert ring into a slice of teleinfo.Frame for JSON marshalling
		mutex.Lock()
		frames := make([]teleinfo.Frame, 0, framesCount)
		framesBuffer.Do(func(v interface{}) {
			if v == nil {
				return
			}
			f := v.(teleinfo.Frame)
			frames = append(frames, f)
		})
		mutex.Unlock()

		// Render JSON
		doc, _ := json.Marshal(frames)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(doc)
	})

	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Starting HTTP service on %s, handling /frames and /metrics\n", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		fmt.Printf("Unable to start HTTP server: %s", err)
		os.Exit(1)
	}
}
