package client

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"

	"sns/app/post/internal/conf"
)

func NewTracerProvider(appInfo *conf.App, traceSetting *conf.Trace) (trace.TracerProvider, error) {
	opts := []tracesdk.TracerProviderOption{
		// set the sampling rate based on the parent span to 100%.
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(traceSetting.SampleRatio))),
		// record information about this application in a resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(appInfo.Name),
			semconv.ServiceVersionKey.String(appInfo.Version),
			semconv.ServiceInstanceIDKey.String(appInfo.Id)),
		),
	}

	if traceSetting.HttpEndpoint != "" {
		// create the jaeger exporter.
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(traceSetting.HttpEndpoint),
			jaeger.WithUsername(traceSetting.Username),
			jaeger.WithPassword(traceSetting.Password),
		))
		if err != nil {
			return nil, err
		}

		// always be sure to batch in production.
		opts = append(opts, tracesdk.WithBatcher(exp))
	}

	tp := tracesdk.NewTracerProvider(opts...)

	return tp, nil
}
