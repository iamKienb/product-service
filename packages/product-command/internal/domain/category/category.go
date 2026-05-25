package category

import (
	"time"

	"product-command-module/internal/domain/shared"

	"github.com/google/uuid"
)

type Category struct {
	ID       shared.CategoryID
	ShopID   uuid.UUID
	Status   string
	ParentID shared.CategoryID
	Name     string
	Slug     string

	CreatedBy uuid.UUID
	UpdatedBy *uuid.UUID

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func New() *Category {
	return nil
}
