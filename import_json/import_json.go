package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

type Record struct {
	HC        uint      `json:"HC_wh"`
	HP        uint      `json:"HP_wh"`
	Power     int       `json:"PAPP_va"`
	Timestamp time.Time `json:"Timestamp"`
	IsHP      bool      `json:"IsHP"`
}

func readRecord(reader *bufio.Reader) (*Record, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, nil
	}
	var e Record
	err = json.Unmarshal(line, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func accumulator(in, out chan *Record) {

	const RecordListCapacity uint = 32
	recordList := make([]*Record, 0, RecordListCapacity)
	var currentMinute int64

	for {
		e := <-in
		if e == nil {
			out <- nil
			return
		}
		// Round to minute
		epoch := e.Timestamp.Unix()
		epoch -= epoch % 60
		if currentMinute == 0 {
			currentMinute = epoch
		}
		glog.V(1).Infof("Accumulating records for minute %d", currentMinute)

		if epoch == currentMinute {
			recordList = append(recordList, e)
			continue
		}

		firstRecord := recordList[0]
		lastRecord := recordList[len(recordList)-1]
		var power int
		for _, r := range recordList {
			power += r.Power
		}
		power = power / len(recordList)
		res := Record{
			Timestamp: time.Unix(currentMinute, 0),
			IsHP:      firstRecord.IsHP,
			HC:        lastRecord.HC - firstRecord.HC,
			HP:        lastRecord.HP - firstRecord.HP,
			Power:     power,
		}
		out <- &res

		currentMinute = epoch
		recordList = make([]*Record, 0, RecordListCapacity)
	}
}

func writePoints(clt influxdb.Client, cfg influxdb.BatchPointsConfig, pts []*influxdb.Point) error {
	bp, err := influxdb.NewBatchPoints(cfg)
	if err != nil {
		glog.Exitf("influxdb.NewBatchPoints failed: %s", err)
	}
	bp.AddPoints(pts)
	return clt.Write(bp)
}

func injector(clt influxdb.Client, cfg influxdb.BatchPointsConfig, in chan *Record, done chan bool) {
	const (
		hcMeasurement = "hc"
		hpMeasurement = "hp"
		batchSize     = 360
	)
	tags := map[string]string{}
	var name string
	var wh uint
	var inserted int
	pts := make([]*influxdb.Point, 0, batchSize)

	for {
		r := <-in
		if r == nil {
			if len(pts) != 0 {
				err := writePoints(clt, cfg, pts)
				if err != nil {
					glog.Errorf("Error adding points to InfluxDB: %s", err)
				} else {
					inserted += len(pts)
				}
			}
			fmt.Printf("Finished inserting %d points\n", inserted)
			done <- true
			return
		}

		if r.IsHP {
			name = hpMeasurement
			wh = r.HP
		} else {
			name = hcMeasurement
			wh = r.HC
		}
		fields := map[string]interface{}{
			"power_va": uint(r.Power),
			"conso_wh": wh,
		}

		pt, err := influxdb.NewPoint(name, tags, fields, r.Timestamp)
		if err != nil {
			glog.Errorf("influxdb.NewPoint failed: %s", err)
			continue
		}
		pts = append(pts, pt)
		if len(pts) == batchSize {
			err = writePoints(clt, cfg, pts)
			if err != nil {
				glog.Errorf("Error adding points to InfluxDB: %s", err)
				continue
			}
			pts = make([]*influxdb.Point, 0, batchSize)
			inserted += batchSize
			fmt.Printf("Inserted %d points, up to %s\n", inserted, r.Timestamp.String())
		}
	}
}

func newClient() (influxdb.Client, error) {
	cfg := influxdb.HTTPConfig{
		Addr:     os.Getenv("INFLUX_URL"),
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PASS"),
		Timeout:  time.Second,
	}
	if len(cfg.Addr) == 0 {
		cfg.Addr = "http://localhost:8086/"
	}
	clt, err := influxdb.NewHTTPClient(cfg)
	if err != nil {
		return nil, err
	}
	_, status, err := clt.Ping(time.Second)
	if err != nil {
		return nil, err
	}
	glog.Infof("Connected to InfluxDB %s", status)
	return clt, nil
}

func readFile(filename string, out chan *Record) {
	file, err := os.Open(filename)
	if err != nil {
		glog.Errorf("Error opening %s: %s", filename, err)
		return
	}
	defer file.Close()

	glog.V(1).Infof("Reading file %s...", filename)
	reader := bufio.NewReader(file)
	for {
		e, err := readRecord(reader)
		if err != nil {
			glog.Errorf("%s: error reading entry (%s)", filename, err)
			continue
		}
		out <- e
		if e == nil {
			return
		}
	}
}

func main() {

	flag.Parse()

	clt, err := newClient()
	if err != nil {
		glog.Exitf("Error connecting to InfluxDB: %s", err)
	}
	defer clt.Close()

	chRecords := make(chan *Record)
	chSamples := make(chan *Record)
	chDone := make(chan bool)
	go accumulator(chRecords, chSamples)
	cfg := influxdb.BatchPointsConfig{
		Precision:       "s",
		Database:        os.Getenv("INFLUX_DB"),
		RetentionPolicy: os.Getenv("INFLUX_RETENTION"),
	}
	go injector(clt, cfg, chSamples, chDone)

	for _, filename := range flag.Args() {
		readFile(filename, chRecords)
	}
	<-chDone
}
