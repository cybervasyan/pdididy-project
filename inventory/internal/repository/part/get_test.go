package part

import (
	"context"
	"testing"
	"time"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	existingPart := model.Part{
		PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Name:          "shitski",
		Description:   "maybe some shitski?",
		Price:         0.15,
		StockQuantity: 67,
		Category:      model.CategoryEngine,
		Dimensions:    &model.Dimensions{Length: 105, Width: 105, Height: 105, Weight: 10},
		Manufacturer:  &model.Manufacturer{Name: "shiiii", Country: "Russia", Website: "shiii.com"},
		Tags:          []string{"shi", "shii", "shiii"},
		Metadata:      map[string]model.Value{"shi?": {Kind: model.ValueKindString, StringValue: "shi"}},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name     string
		stored   []model.Part
		reqUUID  uuid.UUID
		wantPart model.Part
		wantErr  error
	}{
		{
			name:     "success",
			stored:   []model.Part{existingPart},
			reqUUID:  existingPart.PartUUID,
			wantPart: existingPart,
			wantErr:  nil,
		},
		{
			name:     "no such part",
			stored:   []model.Part{existingPart},
			reqUUID:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			wantPart: model.Part{},
			wantErr:  model.ErrPartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewRepository(tt.stored)

			got, err := repo.Get(context.Background(), tt.reqUUID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, model.Part{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPart, got)
		})
	}
}
