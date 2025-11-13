package order

import (
	"context"

	"github.com/Denisz0785/spaceyard/order/internal/model"
	repoModel "github.com/Denisz0785/spaceyard/order/internal/repo/model"
	"github.com/google/uuid"
)

func (s *storage) Get(ctx context.Context, uuid uuid.UUID) (repoModel.Order, error) {
	orderUUID := uuid.String()

	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
