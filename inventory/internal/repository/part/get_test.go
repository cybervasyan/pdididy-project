package part

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

type collectionMock struct {
	findOneResult *mongo.SingleResult
	findCursor    *mongo.Cursor
	findErr       error
	filter        any
}

func (m *collectionMock) FindOne(_ context.Context, filter any, _ ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	m.filter = filter
	return m.findOneResult
}

func (m *collectionMock) Find(_ context.Context, filter any, _ ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	m.filter = filter
	return m.findCursor, m.findErr
}

func TestGet(t *testing.T) {
	t.Parallel()

	partUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	want := model.Part{PartUUID: partUUID, Name: "Main Engine", Price: 100}
	driverErr := errors.New("mongo failure")

	tests := []struct {
		name    string
		result  *mongo.SingleResult
		want    model.Part
		wantErr error
	}{
		{
			name:   "success",
			result: mongo.NewSingleResultFromDocument(want, nil, nil),
			want:   want,
		},
		{
			name:    "not found",
			result:  mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil),
			wantErr: model.ErrPartNotFound,
		},
		{
			name:    "driver error",
			result:  mongo.NewSingleResultFromDocument(bson.D{}, driverErr, nil),
			wantErr: driverErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collection := &collectionMock{findOneResult: tt.result}
			repo := NewRepository(collection)

			got, err := repo.Get(context.Background(), partUUID)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, bson.M{"part_uuid": partUUID}, collection.filter)
		})
	}
}
