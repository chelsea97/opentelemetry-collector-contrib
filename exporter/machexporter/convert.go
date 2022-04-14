package machexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/machexporter"

import (
	"fmt"
	"go.opentelemetry.io/collector/model/pdata"
	pb "github.com/fsolleza/mach-proto/golang"
)

// https://github.com/open-telemetry/opentelemetry-collector/blob/7d2df7b067fed4cd366c4f7502f218b4a4d4a723/pdata/internal/generated_ptrace.go
func handle_traces(td pdata.Traces) {
	fmt.Println("TRACE: *****")
	resource_spans_slice := td.ResourceSpans()
	resource_spans := resource_spans_slice.At(0)
	scope_spans_slice := resource_spans.ScopeSpans()
	scope_spans := scope_spans_slice.At(0)
	span_slice := scope_spans.Spans()
	span := span_slice.At(0)

	trace_id := span.TraceID()
	span_id := span.SpanID()
	parent_span := span.ParentSpanID()
	name := span.Name()
	start_time := span.StartTimestamp()
	end_time := span.EndTimestamp()
	attribs := span.Attributes().AsRaw()
	fmt.Println("TraceID", trace_id)
	fmt.Println("SpanId", span_id)
	fmt.Println("Parent Span", parent_span)
	fmt.Println("Name", name)
	fmt.Println("Start Time", start_time)
	fmt.Println("End Time", end_time)
	fmt.Println("Attrib", attribs)
}

func handle_sum(metric pdata.Metric) {
	var types []pb.AddSeriesRequest_ValueType
	types = append(types, pb.AddSeriesRequest_F64)
	fmt.Println("random print to get pb", types)
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


