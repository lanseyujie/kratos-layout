package client

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"

	"sns/app/post/internal/conf"
)

func NewTracerProvider(appInfo *conf.App, traceSetting *conf.Trace) (trace.TracerProvider, error) {
	opts := make([]tracesdk.TracerProviderOption, 0, 4)

	fraction := 1.0
	if traceSetting.SampleRatio != nil {
		fraction = *traceSetting.SampleRatio
	}

	// set the sampling rate based on the parent span to 100%.
	opts = append(opts, tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(fraction))))

	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(appInfo.Name),
		semconv.ServiceVersionKey.String(appInfo.Version),
		semconv.ServiceInstanceIDKey.String(appInfo.Id),
	}

	// record information about this application in a resource.
	opts = append(opts, tracesdk.WithResource(resource.NewSchemaless(attrs...)))

	if traceSetting.HttpEndpoint != nil && *traceSetting.HttpEndpoint != "" {
		// create the jaeger exporter.
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(*traceSetting.HttpEndpoint)))
		if err != nil {
			return nil, err
		}

		// always be sure to batch in production.
		opts = append(opts, tracesdk.WithBatcher(exp))
	}

	tp := tracesdk.NewTracerProvider(opts...)

	return tp, nil
}
