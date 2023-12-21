package mongotel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedCollection struct {
	CollectionOperations
}

func TraceCollection(coll CollectionOperations) CollectionOperations {
	return tracedCollection{
		CollectionOperations: coll,
	}
}

func (t tracedCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (rs *mongo.SingleResult) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("FindOne", trace.WithAttributes(
			attribute.String("Filter", fmt.Sprintf("%v", filter)),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, rs.Err())
	}(time.Now())

	return t.CollectionOperations.FindOne(ctx, filter, opts...)
}

func (t tracedCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (rs *mongo.InsertOneResult, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("FindOne", trace.WithAttributes(
			attribute.String("Document", fmt.Sprintf("%v", document)),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.CollectionOperations.InsertOne(ctx, document, opts...)
}

func (t tracedCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (rs *mongo.UpdateResult, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("FindOne", trace.WithAttributes(
			attribute.String("Filter", fmt.Sprintf("%v", filter)),
			attribute.String("Update", fmt.Sprintf("%v", update)),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.CollectionOperations.UpdateOne(ctx, filter, update, opts...)
}

func (t tracedCollection) recordError(span trace.Span, err error) {
	if err != nil {
		var cmdErr mongo.CommandError
		if errors.As(err, &cmdErr) {
			span.AddEvent("Command Error", trace.WithAttributes(
				attribute.Int("Code", int(cmdErr.Code)),
				attribute.String("Message", cmdErr.Message),
				attribute.StringSlice("Labels", cmdErr.Labels),
				attribute.String("Name", cmdErr.Name),
				attribute.String("Wrapped", fmt.Sprintf("%v", cmdErr.Wrapped)),
				attribute.String("Raw", fmt.Sprintf("%s", cmdErr.Raw)),
			))
		} else {
			span.AddEvent("Error", trace.WithAttributes(
				attribute.String("Message", err.Error()),
			))
		}
	}
}
