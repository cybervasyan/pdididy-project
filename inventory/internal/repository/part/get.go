package part

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	filter := bson.M{
		"part_uuid": partUUID,
	}

	var part model.Part

	err := r.collection.FindOne(ctx, filter).Decode(&part)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Part{}, model.ErrPartNotFound
	}

	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
