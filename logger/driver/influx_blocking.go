package driver

import (
	"context"
	"fmt"
	"time"

	"github.com/hugebear-io/gofiber/fabric"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap/zapcore"
)

type influxWriterBlocking struct {
	measurement string
	writeAPI    api.WriteAPIBlocking
	ctx         context.Context
}

func NewInfluxWriterBlocking(client influxdb2.Client, org, bucket, measurement string) *zapcore.WriteSyncer {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	writer := &influxWriterBlocking{writeAPI: writeAPI, ctx: context.Background(), measurement: measurement}
	syncer := zapcore.AddSync(writer)
	return &syncer
}

func (w influxWriterBlocking) Write(p []byte) (int, error) {
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
	w.writeAPI.WritePoint(context.Background(), point)
	return len(p), nil
}
