package part

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	def "github.com/cybervasyan/pdididy-project/inventory/internal/repository"
)

var _ def.Repository = (*repository)(nil)

type repository struct {
	collection collection
}

type collection interface {
	FindOne(context.Context, any, ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	Find(context.Context, any, ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
}

func NewRepository(collection collection) *repository {
	return &repository{
		collection: collection,
	}
}
