package order

import (
	"context"
	"github.com/Denisz0785/spaceyard/order/internal/repo/converter"

	"github.com/Denisz0785/spaceyard/order/internal/model"
)

func (s *storage) Update(_ context.Context, order *model.Order) error {
	orderUUID := order.OrderUUID.String()

	s.mu.Lock() // Lock for the entire read-modify-write operation
	defer s.mu.Unlock()

	_, ok := s.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	s.orders[orderUUID] = converter.OrderToRepoOrder(order)

	return nil
}
