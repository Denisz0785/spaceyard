package v1

import (
	"context"
	"github.com/Denisz0785/spaceyard/order/internal/converter"
	orderv1 "github.com/Denisz0785/spaceyard/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	resp, err := a.orderService.CreateOrder(ctx, converter.OrderInfoToModel(req))
	if err != nil {
		return nil, err
	}

	return converter.ModelToCreateOrderResponse(*resp), nil
}
