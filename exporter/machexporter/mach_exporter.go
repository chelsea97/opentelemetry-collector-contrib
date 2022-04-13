// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/machexporter"

import (
	"context"
	"sync"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/otlp"
	"go.opentelemetry.io/collector/model/pdata"
)

// Marshaler configuration used for marhsaling Protobuf to JSON.
var tracesMarshaler = otlp.NewJSONTracesMarshaler()
var metricsMarshaler = otlp.NewJSONMetricsMarshaler()
var logsMarshaler = otlp.NewJSONLogsMarshaler()

// fileExporter is the implementation of file exporter that writes telemetry data to a file
// in Protobuf-JSON format.
type machExporter struct {
	mutex sync.Mutex
}

func (e *machExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *machExporter) ConsumeTraces(_ context.Context, td pdata.Traces) error {
	buf, err := tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func handle_sum(metric pdata.Metric) {
	point_slice := metric.Sum().DataPoints()
	point := point_slice.At(0)
	timestamp := point.Timestamp()
	attribs := point.Attributes().AsRaw()
	value_type := point.ValueType()
	fmt.Println(timestamp)
	fmt.Println(value_type)
	fmt.Println(attribs)
	switch value_type {
	case pdata.MetricValueTypeInt:
		fmt.Println(point.IntVal())
	case pdata.MetricValueTypeDouble:
		fmt.Println(point.DoubleVal())
	default:
		fmt.Println("FAILED")
	}
	fmt.Println("SUM HERE")
}

func handle_histogram(metric pdata.Metric) {
	point_slice := metric.Histogram().DataPoints()
	point := point_slice.At(0)
	timestamp := point.Timestamp()
	attribs := point.Attributes().AsRaw()

	count := point.Count()
	sum := 0.0
	if point.HasSum() {
		sum = point.Sum()
	}
	bucket_counts:= point.BucketCounts()
	bounds:= point.ExplicitBounds()

	fmt.Println("Timestamp", timestamp)
	fmt.Println("Attribs", attribs)
	fmt.Println("Count", count)
	fmt.Println("Sum", sum)
	fmt.Println("bucket counts", bucket_counts)
	fmt.Println("bounds", bounds)
	fmt.Println("HISTOGRAM HERE")
}

// This stuff can be found here:
// https://github.com/open-telemetry/opentelemetry-collector/blob/7d2df7b067fed4cd366c4f7502f218b4a4d4a723/pdata/internal/generated_pmetric.go
// 
func handle_metrics(md pdata.Metrics) {
	fmt.Println("*************")
	resource_slice := md.ResourceMetrics()
	resource := resource_slice.At(0)
	scope_slice := resource.ScopeMetrics()
	scope := scope_slice.At(0)
	fmt.Println("Scope name", scope.Scope().Name())
	metrics_slice := scope.Metrics()
	metric := metrics_slice.At(0)
	fmt.Println("Metric name", metric.Name())
	datatype := metric.DataType()
	fmt.Println(datatype)

	//fmt.Println(md.ResourceMetrics().Len())

	switch datatype {
	case pdata.MetricDataTypeSum:
		handle_sum(metric)
		//fmt.Println(points.Value)
	case pdata.MetricDataTypeGauge:
		fmt.Println("GAUGE HERE")
	case pdata.MetricDataTypeHistogram:
		handle_histogram(metric)
	case pdata.MetricDataTypeExponentialHistogram:
		fmt.Println("EXP HIST HERE")
	case pdata.MetricDataTypeSummary:
		fmt.Println("SUMMARY HERE")
	}
}

func (e *machExporter) ConsumeMetrics(_ context.Context, md pdata.Metrics) error {
	handle_metrics(md)
	buf, err := metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func (e *machExporter) ConsumeLogs(_ context.Context, ld pdata.Logs) error {
	buf, err := logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return exportMessageAsLine(e, buf)
}

func exportMessageAsLine(e *machExporter, buf []byte) error {
	// Ensure only one write operation happens at a time.
	e.mutex.Lock()
	defer e.mutex.Unlock()
	var toPrint string
	toPrint = string(buf)
	fmt.Println(toPrint)
	return nil
}

func (e *machExporter) Start(context.Context, component.Host) error {
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *machExporter) Shutdown(context.Context) error {
	return nil
}
