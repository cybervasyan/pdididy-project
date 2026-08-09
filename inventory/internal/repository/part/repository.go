package part

import (
	"go.mongodb.org/mongo-driver/v2/mongo"

	def "github.com/cybervasyan/pdididy-project/inventory/internal/repository"
)

var _ def.Repository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *repository {
	return &repository{
		collection: collection,
	}
}
