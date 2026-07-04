package part

import (
	"errors"
	"time"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (s *ServiceSuite) TestList() {
	fixedTime := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	unexpectedErr := errors.New("repository failed")

	repoParts := []repoModel.Part{
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
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
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "Fuel Pressure Valve",
			Description:   "Valve for stabilizing fuel line pressure",
			Price:         45.50,
			StockQuantity: 34,
			Category:      repoModel.CategoryFuel,
			Dimensions:    &repoModel.Dimensions{Length: 12, Width: 4.5, Height: 4.5, Weight: 0.8},
			Manufacturer:  &repoModel.Manufacturer{Name: "FuelTech Labs", Country: "Japan", Website: "fueltech.example"},
			Tags:          []string{"fuel", "valve", "pressure"},
			Metadata: map[string]repoModel.Value{
				"max_pressure_bar": {Kind: repoModel.ValueKindDouble, DoubleValue: 18.5},
				"certified":        {Kind: repoModel.ValueKindBool, BoolValue: true},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:          "Reinforced Porthole Glass",
			Description:   "Impact-resistant glass panel for cabin porthole",
			Price:         310.75,
			StockQuantity: 6,
			Category:      repoModel.CategoryPorthole,
			Dimensions:    &repoModel.Dimensions{Length: 54, Width: 54, Height: 3.2, Weight: 7.6},
			Manufacturer:  &repoModel.Manufacturer{Name: "ClearSky Components", Country: "France", Website: "clearsky.example"},
			Tags:          []string{"porthole", "glass", "cabin"},
			Metadata: map[string]repoModel.Value{
				"glass_type": {Kind: repoModel.ValueKindString, StringValue: "laminated"},
				"layers":     {Kind: repoModel.ValueKindInt64, Int64Value: 4},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:          "Left Wing Stabilizer",
			Description:   "Stabilizer module for left wing assembly",
			Price:         780.00,
			StockQuantity: 3,
			Category:      repoModel.CategoryWing,
			Dimensions:    &repoModel.Dimensions{Length: 180, Width: 42, Height: 16, Weight: 24.3},
			Manufacturer:  &repoModel.Manufacturer{Name: "AeroFrame Works", Country: "USA", Website: "aeroframe.example"},
			Tags:          []string{"wing", "stabilizer", "aero"},
			Metadata: map[string]repoModel.Value{
				"side":                 {Kind: repoModel.ValueKindString, StringValue: "left"},
				"requires_calibration": {Kind: repoModel.ValueKindBool, BoolValue: true},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:          "Generic Inspection Plug",
			Description:   "Temporary plug for diagnostics and maintenance",
			Price:         5.25,
			StockQuantity: 150,
			Category:      repoModel.CategoryUnspecified,
			Dimensions:    &repoModel.Dimensions{Length: 4, Width: 4, Height: 2, Weight: 0.1},
			Manufacturer:  &repoModel.Manufacturer{Name: "Service Parts Co", Country: "Russia", Website: "service-parts.example"},
			Tags:          []string{"maintenance", "plug", "inspection"},
			Metadata: map[string]repoModel.Value{
				"disposable": {Kind: repoModel.ValueKindBool, BoolValue: true},
				"batch":      {Kind: repoModel.ValueKindString, StringValue: "A-2026-07"},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}
	wantParts := []model.Part{
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
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
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "Fuel Pressure Valve",
			Description:   "Valve for stabilizing fuel line pressure",
			Price:         45.50,
			StockQuantity: 34,
			Category:      model.CategoryFuel,
			Dimensions:    &model.Dimensions{Length: 12, Width: 4.5, Height: 4.5, Weight: 0.8},
			Manufacturer:  &model.Manufacturer{Name: "FuelTech Labs", Country: "Japan", Website: "fueltech.example"},
			Tags:          []string{"fuel", "valve", "pressure"},
			Metadata: map[string]model.Value{
				"max_pressure_bar": {Kind: model.ValueKindDouble, DoubleValue: 18.5},
				"certified":        {Kind: model.ValueKindBool, BoolValue: true},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:          "Reinforced Porthole Glass",
			Description:   "Impact-resistant glass panel for cabin porthole",
			Price:         310.75,
			StockQuantity: 6,
			Category:      model.CategoryPorthole,
			Dimensions:    &model.Dimensions{Length: 54, Width: 54, Height: 3.2, Weight: 7.6},
			Manufacturer:  &model.Manufacturer{Name: "ClearSky Components", Country: "France", Website: "clearsky.example"},
			Tags:          []string{"porthole", "glass", "cabin"},
			Metadata: map[string]model.Value{
				"glass_type": {Kind: model.ValueKindString, StringValue: "laminated"},
				"layers":     {Kind: model.ValueKindInt64, Int64Value: 4},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:          "Left Wing Stabilizer",
			Description:   "Stabilizer module for left wing assembly",
			Price:         780.00,
			StockQuantity: 3,
			Category:      model.CategoryWing,
			Dimensions:    &model.Dimensions{Length: 180, Width: 42, Height: 16, Weight: 24.3},
			Manufacturer:  &model.Manufacturer{Name: "AeroFrame Works", Country: "USA", Website: "aeroframe.example"},
			Tags:          []string{"wing", "stabilizer", "aero"},
			Metadata: map[string]model.Value{
				"side":                 {Kind: model.ValueKindString, StringValue: "left"},
				"requires_calibration": {Kind: model.ValueKindBool, BoolValue: true},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:          "Generic Inspection Plug",
			Description:   "Temporary plug for diagnostics and maintenance",
			Price:         5.25,
			StockQuantity: 150,
			Category:      model.CategoryUnspecified,
			Dimensions:    &model.Dimensions{Length: 4, Width: 4, Height: 2, Weight: 0.1},
			Manufacturer:  &model.Manufacturer{Name: "Service Parts Co", Country: "Russia", Website: "service-parts.example"},
			Tags:          []string{"maintenance", "plug", "inspection"},
			Metadata: map[string]model.Value{
				"disposable": {Kind: model.ValueKindBool, BoolValue: true},
				"batch":      {Kind: model.ValueKindString, StringValue: "A-2026-07"},
			},
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	tests := []struct {
		name       string
		reqService model.PartsFilter
		reqRepo    repoModel.PartsFilter
		repoParts  []repoModel.Part
		repoErr    error
		wantParts  []model.Part
		wantErr    error
	}{
		{
			name:       "success",
			reqService: model.PartsFilter{},
			reqRepo: repoModel.PartsFilter{
				Categories: []repoModel.Category{},
			},
			repoParts: repoParts,
			repoErr:   nil,
			wantParts: wantParts,
			wantErr:   nil,
		},
		{
			name: "filter by uuid",
			reqService: model.PartsFilter{
				PartUUIDs: []uuid.UUID{wantParts[0].PartUUID},
			},
			reqRepo: repoModel.PartsFilter{
				PartUUIDs:  []uuid.UUID{repoParts[0].PartUUID},
				Categories: []repoModel.Category{},
			},
			repoParts: []repoModel.Part{repoParts[0]},
			wantParts: []model.Part{wantParts[0]},
			wantErr:   nil,
		},
		{
			name: "filter by category",
			reqService: model.PartsFilter{
				Categories: []model.Category{model.CategoryFuel},
			},
			reqRepo: repoModel.PartsFilter{
				Categories: []repoModel.Category{repoModel.CategoryFuel},
			},
			repoParts: []repoModel.Part{repoParts[1]},
			wantParts: []model.Part{wantParts[1]},
			wantErr:   nil,
		},
		{
			name: "filter by manufacturer country",
			reqService: model.PartsFilter{
				ManufacturerCountries: []string{"Russia"},
			},
			reqRepo: repoModel.PartsFilter{
				ManufacturerCountries: []string{"Russia"},
				Categories:            []repoModel.Category{},
			},
			repoParts: []repoModel.Part{repoParts[4]},
			wantParts: []model.Part{wantParts[4]},
			wantErr:   nil,
		},
		{
			name: "filter by tag",
			reqService: model.PartsFilter{
				Tags: []string{"fuel"},
			},
			reqRepo: repoModel.PartsFilter{
				Tags:       []string{"fuel"},
				Categories: []repoModel.Category{},
			},
			repoParts: []repoModel.Part{repoParts[0], repoParts[1]},
			wantParts: []model.Part{wantParts[0], wantParts[1]},
			wantErr:   nil,
		},
		{
			name: "filter by several fields",
			reqService: model.PartsFilter{
				Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
				ManufacturerCountries: []string{"Japan"},
				Tags:                  []string{"pressure"},
			},
			reqRepo: repoModel.PartsFilter{
				Categories:            []repoModel.Category{repoModel.CategoryEngine, repoModel.CategoryFuel},
				ManufacturerCountries: []string{"Japan"},
				Tags:                  []string{"pressure"},
			},
			repoParts: []repoModel.Part{repoParts[1]},
			wantParts: []model.Part{wantParts[1]},
			wantErr:   nil,
		},
		{
			name: "empty result",
			reqService: model.PartsFilter{
				Names: []string{"Missing Part"},
			},
			reqRepo: repoModel.PartsFilter{
				Names:      []string{"Missing Part"},
				Categories: []repoModel.Category{},
			},
			repoParts: []repoModel.Part{},
			wantParts: []model.Part{},
			wantErr:   nil,
		},
		{
			name:       "repository error",
			reqService: model.PartsFilter{},
			reqRepo: repoModel.PartsFilter{
				Categories: []repoModel.Category{},
			},
			repoParts: []repoModel.Part{},
			repoErr:   unexpectedErr,
			wantParts: []model.Part{},
			wantErr:   unexpectedErr,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryRepo.EXPECT().
				List(s.ctx, tt.reqRepo).
				Return(tt.repoParts, tt.repoErr).
				Once()

			got, err := s.service.ListParts(s.ctx, tt.reqService)

			if tt.wantErr != nil {
				require.ErrorIs(s.T(), err, tt.wantErr)
				require.Empty(s.T(), got)
				return
			}

			require.NoError(s.T(), err)
			require.ElementsMatch(s.T(), tt.wantParts, got)
		})
	}
}
