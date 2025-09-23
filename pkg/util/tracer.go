package util

import (
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func RecordSpanError(span trace.Span, err error, f string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Do nothing
		return
	}
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err, oteltrace.WithAttributes(attribute.String("func", f)))
}
