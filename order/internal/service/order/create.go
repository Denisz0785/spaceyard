package order

import (
	"context"
	"fmt"

	"github.com/Denisz0785/spaceyard/order/internal/model"
)

func (s *orderService) CreateOrder(ctx context.Context, orderInfo *model.CreateOrderInfo) (*model.CreateOrderResponse, error) {
	// 1. Call InventoryService to get part details
	partUUIDsStrings := make([]string, len(orderInfo.PartUuids))
	for i, u := range orderInfo.PartUuids {
		partUUIDsStrings[i] = u.String()
	}

	inventoryResp, err := s.inventoryClient.ListParts(ctx, partUUIDsStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to communicate with inventory service: %w", err)
	}

	// 2. Validate that all requested parts were found
	if len(inventoryResp) != len(orderInfo.PartUuids) {
		return nil, fmt.Errorf("one or more requested parts do not exist")
	}

	// 3. Calculate total price
	var totalPrice float64
	for _, part := range inventoryResp {
		totalPrice += part.Price
	}

	order := &model.Order{
		UserUUID:        orderInfo.UserUUID,
		PartUuids:       orderInfo.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPENDINGPAYMENT,
	}

	uuid, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	resp := &model.CreateOrderResponse{
		OrderUUID:  uuid,
		TotalPrice: order.TotalPrice,
	}

	return resp, nil
}
