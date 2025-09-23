package repositories

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/app"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/pkg/cache"
	"github.com/tomatosAt/reskill-go-react/pkg/database"
	"github.com/tomatosAt/reskill-go-react/pkg/requests"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

const moduleName = "openapi"

type Repository struct {
	app     *app.Context
	http    *requests.HttpClient
	log     *logrus.Entry
	tracer  trace.Tracer
	dbMain  *database.Client
	dbMongo *mongo.Database
	dbLog   *database.Client
	cache   *cache.Redis
}

func (r Repository) Module() string {
	return moduleName
}

func (r Repository) AppCfg() *config.Config {
	return r.app.Config
}

func (r Repository) Log() *logrus.Entry {
	return r.log.Dup()
}

func (r Repository) DB() *database.Client {
	return r.dbMain
}

func (r Repository) Cache() *cache.Redis {
	return r.cache
}

func (r Repository) DBMongo() *mongo.Database {
	return r.dbMongo
}

func (r Repository) Trace(ctx context.Context, spanName string, attributes ...trace.SpanStartOption) (context.Context, trace.Span) {
	return r.tracer.Start(ctx, spanName, attributes...)
}

// func (r Repository) Discord() *discord.Discord {
// 	return r.discord
// }

func New(app *app.Context) (*Repository, error) {
	return &Repository{}, nil
}
