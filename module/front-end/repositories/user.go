package repositories

import (
	"context"
	"errors"

	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func (r *Repository) GetAuthCtxRepo(ctx context.Context) (*dto.AuthSession, error) {
	_, span := r.Trace(ctx, "repo.GetAuthCtxRepo", oteltrace.WithAttributes())
	defer span.End()
	authSession, ok := ctx.Value(config.AuthorizeContext).(dto.AuthSession)
	if !ok {
		err := errors.New("unauthorized claims")
		util.RecordSpanError(span, err, "repo.GetAuthCtxRepo")
		return nil, errors.New("unauthorized claims")
	}
	return &authSession, nil
}
