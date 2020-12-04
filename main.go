package main

import (
	"bytes"
	"flag"
	dto "github.com/prometheus/client_model/go"
	. "github.com/prometheus/common/expfmt"
	"github.com/shirou/gopsutil/cpu"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func printPrometheus(w http.ResponseWriter, r *http.Request) {
	c,_ := cpu.Times(false)

	var buff bytes.Buffer
	delimEncoder := NewEncoder(&buff, FmtText)
	metric := &dto.MetricFamily{
		Name: proto.String("cpu_info"),
		Type: dto.MetricType_GAUGE.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair {
					{
						Name: proto.String("type"),
						Value: proto.String("user"),
					},
				},
				Gauge: &dto.Gauge{
					Value: proto.Float64(c[0].User),
				},
			},
		},
	}

	delimEncoder.Encode(metric)

	metric = &dto.MetricFamily{
		Name: proto.String("cpu_info"),
		Type: dto.MetricType_GAUGE.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair {
					{
						Name: proto.String("type"),
						Value: proto.String("system"),
					},
				},
				Gauge: &dto.Gauge{
					Value: proto.Float64(c[0].System),
				},
			},
		},
	}
	delimEncoder.Encode(metric)

	w.Write(buff.Bytes())
}

func main() {
	flag.Parse()
	http.HandleFunc("/metrics",printPrometheus)
	//http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}