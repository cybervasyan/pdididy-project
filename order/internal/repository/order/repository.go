package order

import (
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

type uuidArray []uuid.UUID

func (ids uuidArray) Value() (driver.Value, error) {
	values := make([]string, len(ids))
	for i, id := range ids {
		values[i] = id.String()
	}

	return "{" + strings.Join(values, ",") + "}", nil
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
