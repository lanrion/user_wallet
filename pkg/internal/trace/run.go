package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func DoWithSpan(ctx context.Context, operationName string, doFn func(ctx context.Context, span trace.Span),
	opts ...trace.SpanStartOption) {
	tracer := otel.Tracer("user-wallet")
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()
	doFn(newCtx, span)
}
