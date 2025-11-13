package v1

import (
	"context"
	"errors"

	"github.com/Denisz0785/spaceyard/order/internal/model"

	orderv1 "github.com/Denisz0785/spaceyard/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrUpdateOrder) {
			return &orderv1.CancelOrderConflict{}, nil
		} else if errors.Is(err, model.ErrOrderNotFound) {
			return &orderv1.CancelOrderNotFound{}, nil
		} else if errors.Is(err, model.ErrCancelOrder) {
			return &orderv1.CancelOrderConflict{}, nil
		}
	}

	return &orderv1.CancelOrderNoContent{}, nil
}
