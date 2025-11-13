package repo

import (
	"context"
	"github.com/google/uuid"

	"github.com/Denisz0785/spaceyard/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (uuid.UUID, error)
	Get(ctx context.Context, uuuid uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
