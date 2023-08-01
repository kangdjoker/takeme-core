package utils

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	// opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func LoggingAPICall(statusCode int, request interface{}, response interface{}, message string) {
	a, _ := json.Marshal(request)
	log.Info("Request body to "+message, string(a))

	a, _ = json.Marshal(response)
	log.Info("Response body from "+message+" status code "+strconv.Itoa(statusCode)+" : ", string(a))
}
func SetupTracer(serviceName string, spanName string) (*io.Closer, *opentracing.Span) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 1000,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	tracer, trCloser, _ := cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	opentracing.SetGlobalTracer(tracer)
	span, _ := opentracing.StartSpanFromContext(context.Background(), spanName)
	return &trCloser, &span
}
