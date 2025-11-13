package order

import (
	"context"

	"github.com/Denisz0785/spaceyard/order/internal/model"
	"github.com/google/uuid"
)

func (s *orderService) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {

	order, err := s.repo.Get(ctx, orderUUID)
	if err != nil {
		return model.ErrOrderNotFound
	}

	// Contract: If an order is already PAID, it cannot be cancelled.
	if order.Status == model.OrderStatusPAID {
		return model.ErrCancelOrder
	}

	// If the order is PENDING_PAYMENT, update its status to CANCELLED.
	order.Status = model.OrderStatusCANCELLED

	err = s.repo.Update(ctx, &order)
	if err != nil {
		return model.ErrUpdateOrder
	}

	return nil
}
