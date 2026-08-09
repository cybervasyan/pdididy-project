package part

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

func TestList(t *testing.T) {
	t.Parallel()

	partUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	want := []model.Part{{PartUUID: partUUID, Name: "Fuel Tank", Tags: []string{"fuel"}}}
	driverErr := errors.New("mongo failure")

	t.Run("all filters", func(t *testing.T) {
		t.Parallel()

		cursor, err := mongo.NewCursorFromDocuments([]any{want[0]}, nil, nil)
		require.NoError(t, err)

		collection := &collectionMock{findCursor: cursor}
		repo := NewRepository(collection)
		filter := model.PartsFilter{
			PartUUIDs:             []uuid.UUID{partUUID},
			Names:                 []string{"Fuel Tank"},
			Categories:            []model.Category{model.CategoryFuel},
			ManufacturerCountries: []string{"Germany"},
			Tags:                  []string{"fuel"},
		}

		got, err := repo.List(context.Background(), filter)

		require.NoError(t, err)
		require.Equal(t, want, got)
		require.Equal(t, bson.M{
			"part_uuid":            bson.M{"$in": filter.PartUUIDs},
			"name":                 bson.M{"$in": filter.Names},
			"category":             bson.M{"$in": filter.Categories},
			"manufacturer.country": bson.M{"$in": filter.ManufacturerCountries},
			"tags":                 bson.M{"$in": filter.Tags},
		}, collection.filter)
	})

	t.Run("empty filter", func(t *testing.T) {
		t.Parallel()

		cursor, err := mongo.NewCursorFromDocuments(nil, nil, nil)
		require.NoError(t, err)
		collection := &collectionMock{findCursor: cursor}

		got, err := NewRepository(collection).List(context.Background(), model.PartsFilter{})

		require.NoError(t, err)
		require.Empty(t, got)
		require.Equal(t, bson.M{}, collection.filter)
	})

	t.Run("find error", func(t *testing.T) {
		t.Parallel()

		collection := &collectionMock{findErr: driverErr}

		got, err := NewRepository(collection).List(context.Background(), model.PartsFilter{})

		require.ErrorIs(t, err, driverErr)
		require.Nil(t, got)
	})

	t.Run("decode error", func(t *testing.T) {
		t.Parallel()

		cursor, err := mongo.NewCursorFromDocuments([]any{bson.M{"price": "invalid"}}, nil, nil)
		require.NoError(t, err)
		collection := &collectionMock{findCursor: cursor}

		got, err := NewRepository(collection).List(context.Background(), model.PartsFilter{})

		require.Error(t, err)
		require.Nil(t, got)
	})
}
