package part

import (
	"errors"
	"time"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (s *ServiceSuite) TestGet() {
	fixedTime := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	partUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	missingUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	unexpectedErr := errors.New("repository failed")

	repoPart := repoModel.Part{
		PartUUID:      partUUID,
		Name:          "Plasma Injector",
		Description:   "High-pressure injector for engine fuel delivery",
		Price:         129.99,
		StockQuantity: 12,
		Category:      repoModel.CategoryEngine,
		Dimensions:    &repoModel.Dimensions{Length: 32.5, Width: 8.2, Height: 8.2, Weight: 2.4},
		Manufacturer:  &repoModel.Manufacturer{Name: "Orbital Engines", Country: "Germany", Website: "orbital-engines.example"},
		Tags:          []string{"engine", "injector", "fuel"},
		Metadata: map[string]repoModel.Value{
			"material": {Kind: repoModel.ValueKindString, StringValue: "titanium"},
			"revision": {Kind: repoModel.ValueKindInt64, Int64Value: 3},
		},
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	wantPart := model.Part{
		PartUUID:      partUUID,
		Name:          "Plasma Injector",
		Description:   "High-pressure injector for engine fuel delivery",
		Price:         129.99,
		StockQuantity: 12,
		Category:      model.CategoryEngine,
		Dimensions:    &model.Dimensions{Length: 32.5, Width: 8.2, Height: 8.2, Weight: 2.4},
		Manufacturer:  &model.Manufacturer{Name: "Orbital Engines", Country: "Germany", Website: "orbital-engines.example"},
		Tags:          []string{"engine", "injector", "fuel"},
		Metadata: map[string]model.Value{
			"material": {Kind: model.ValueKindString, StringValue: "titanium"},
			"revision": {Kind: model.ValueKindInt64, Int64Value: 3},
		},
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	tests := []struct {
		name     string
		reqUUID  uuid.UUID
		repoPart repoModel.Part
		repoErr  error
		wantPart model.Part
		wantErr  error
	}{
		{
			name:     "success",
			reqUUID:  partUUID,
			repoPart: repoPart,
			repoErr:  nil,
			wantPart: wantPart,
			wantErr:  nil,
		},
		{
			name:     "not found",
			reqUUID:  missingUUID,
			repoPart: repoModel.Part{},
			repoErr:  repoModel.ErrPartNotFound,
			wantPart: model.Part{},
			wantErr:  model.ErrPartNotFound,
		},
		{
			name:     "repository error",
			reqUUID:  partUUID,
			repoPart: repoModel.Part{},
			repoErr:  unexpectedErr,
			wantPart: model.Part{},
			wantErr:  unexpectedErr,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryRepo.EXPECT().
				Get(s.ctx, tt.reqUUID).
				Return(tt.repoPart, tt.repoErr).
				Once()

			got, err := s.service.GetPart(s.ctx, tt.reqUUID)

			if tt.wantErr != nil {
				require.ErrorIs(s.T(), err, tt.wantErr)
				require.Equal(s.T(), tt.wantPart, got)
				return
			}

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.wantPart, got)
		})
	}
}
