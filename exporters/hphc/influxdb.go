package hphc

import (
	"flag"
	"fmt"
	"goconso/exporters"
	"goconso/teleinfo"
	"os"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

const (
	exporterName = "hphc.influxdb"
)

var (
	dbURL       = *flag.String(exporterName+".url", getEnvDefault("INFLUXDB_URL", "http://localhost:8086"), "URL to connect to InfluxDB. Can be set from env var 'INFLUXDB_URL'")
	dbName      = *flag.String(exporterName+".name", getEnvDefault("INFLUXDB_NAME", "goconso"), "Name of the target InfluxDB database. Can be set from env var 'INFLUXDB_NAME'.")
	dbUser      = *flag.String(exporterName+".user", os.Getenv("INFLUXDB_USER"), "Username for InfluxDB authentication. Can be set from env var 'INFLUXDB_USER'.")
	dbPass      = *flag.String(exporterName+".pass", os.Getenv("INFLUXDB_PASS"), "Password for InfluxDB authentication. Can be set from env var 'INFLUXDB_PASS'.")
	dbRetention = *flag.String(exporterName+".retpolicy", getEnvDefault("INFLUXDB_RETPOLICY", "default"), "InfluxDB retention policy target name. Can be set from env var 'INFLUXDB_RETPOLICY'.")
)

func getEnvDefault(name string, fallback string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		return fallback
	}
	return val
}

func init() {
	exporters.Register(exporterName, newInfluxDbExporter)
}

type influxDbExporter struct {
	client   influxdb.Client
	batchCfg influxdb.BatchPointsConfig
}

func newInfluxDbExporter() (teleinfo.Exporter, error) {

	cfg := influxdb.HTTPConfig{
		Addr:     dbURL,
		Username: dbUser,
		Password: dbPass,
		Timeout:  time.Second,
	}
	clt, err := influxdb.NewHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	_, _, err = clt.Ping(0)
	if err != nil {
		return nil, err
	}

	batchCfg := influxdb.BatchPointsConfig{
		Precision:       "s",
		Database:        dbName,
		RetentionPolicy: dbRetention,
	}

	exp := &influxDbExporter{
		client:   clt,
		batchCfg: batchCfg,
	}
	return exp, nil
}

func convertRecord(r *record) map[string]interface{} {
	var wh uint32

	if r.IsHP {
		wh = r.HP_wh
	} else {
		wh = r.HC_wh
	}
	fields := map[string]interface{}{
		"power_va":     r.PAPP_va,
		"abs_conso_wh": wh,
	}
	return fields
}

func targetMeasurement(r *record) string {
	if r.IsHP {
		return "hp_raw"
	}
	return "hc_raw"
}

func (x *influxDbExporter) ExportFrame(f teleinfo.Frame) error {
	r := newRecord(f)
	tags := map[string]string{}

	fields := convertRecord(r)
	measurement := targetMeasurement(r)

	pt, err := influxdb.NewPoint(measurement, tags, fields, r.Timestamp)
	if err != nil {
		return fmt.Errorf("influxdb.NewPoint() failed (%s)", err)
	}

	pts, err := influxdb.NewBatchPoints(x.batchCfg)
	if err != nil {
		return fmt.Errorf("influxdb.NewBatchPoints() failed (%s)", err)
	}
	pts.AddPoint(pt)

	return x.client.Write(pts)
}
