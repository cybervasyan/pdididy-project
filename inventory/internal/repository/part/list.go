package part

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

func (r *repository) List(
	ctx context.Context,
	req model.PartsFilter,
) ([]model.Part, error) {
	filter := bson.M{}

	if len(req.PartUUIDs) > 0 {
		filter["part_uuid"] = bson.M{
			"$in": req.PartUUIDs,
		}
	}

	if len(req.Names) > 0 {
		filter["name"] = bson.M{
			"$in": req.Names,
		}
	}

	if len(req.Categories) > 0 {
		filter["category"] = bson.M{
			"$in": req.Categories,
		}
	}

	if len(req.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{
			"$in": req.ManufacturerCountries,
		}
	}

	if len(req.Tags) > 0 {
		filter["tags"] = bson.M{
			"$in": req.Tags,
		}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	parts := make([]model.Part, 0)

	if err = cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	return parts, nil
}
