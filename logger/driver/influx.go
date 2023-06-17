package driver

import (
	"fmt"
	"time"

	"github.com/hugebear-io/gofiber/fabric"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap/zapcore"
)

type influxWriter struct {
	measurement string
	writeAPI    api.WriteAPI
}

func NewInfluxWriter(client influxdb2.Client, org, bucket, measurement string) *zapcore.WriteSyncer {
	writeAPI := client.WriteAPI(org, bucket)
	writer := &influxWriter{writeAPI: writeAPI, measurement: measurement}
	syncer := zapcore.AddSync(writer)
	return &syncer
}

func (w influxWriter) Write(p []byte) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC] influx writer blocking: %v", r)
		}
	}()

	log := map[string]interface{}{}
	fabric.Recast(p, &log)

	fields := map[string]interface{}{"value": string(p)}
	tagFields := []string{"level", "ip", "type", "method", "path", "status-code", "uid"}
	tags := map[string]string{}
	for _, key := range tagFields {
		if val, ok := log[key]; ok {
			tags[key] = fmt.Sprint(val)
		} else if key == "type" {
			tags[key] = "system"
		} else {
			tags[key] = "-"
		}
	}

	point := write.NewPoint(w.measurement, tags, fields, time.Now())
	w.writeAPI.WritePoint(point)
	return len(p), nil
}
