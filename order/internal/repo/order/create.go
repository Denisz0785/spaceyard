package order

import (
	"context"
	"github.com/Denisz0785/spaceyard/order/internal/repo/converter"
	"github.com/google/uuid"

	"github.com/Denisz0785/spaceyard/order/internal/model"
)

func (s *storage) Create(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	order.OrderUUID = uuid.New()

	s.mu.Lock()
	defer s.mu.Unlock()

	repoOrder := converter.OrderToRepoOrder(order)

	s.orders[order.OrderUUID.String()] = repoOrder

	return order.OrderUUID, nil
}
