package model

import (
	"time"

	"github.com/google/uuid"
)

type PartsFilter struct {
	PartUUIDs             []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Part struct {
	PartUUID      uuid.UUID        `bson:"part_uuid"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float64          `bson:"price"`
	StockQuantity int64            `bson:"stock_quantity"`
	Category      Category         `bson:"category"`
	Dimensions    *Dimensions      `bson:"dimensions"`
	Manufacturer  *Manufacturer    `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at"`
}

type Category string

const (
	CategoryUnspecified Category = "CATEGORY_UNSPECIFIED"
	CategoryEngine      Category = "CATEGORY_ENGINE"
	CategoryFuel        Category = "CATEGORY_FUEL"
	CategoryPorthole    Category = "CATEGORY_PORTHOLE"
	CategoryWing        Category = "CATEGORY_WING"
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

// ValueKind говорит, какое из полей Value валидно (эмуляция proto oneof).
type ValueKind string

const (
	ValueKindUnspecified ValueKind = ""
	ValueKindString      ValueKind = "STRING"
	ValueKindInt64       ValueKind = "INT64"
	ValueKindDouble      ValueKind = "DOUBLE"
	ValueKindBool        ValueKind = "BOOL"
)

type Value struct {
	Kind        ValueKind `bson:"kind"`
	StringValue string    `bson:"string_value"`
	Int64Value  int64     `bson:"int64_value"`
	DoubleValue float64   `bson:"double_value"`
	BoolValue   bool      `bson:"bool_value"`
}
