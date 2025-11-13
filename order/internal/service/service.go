package service

import (
	"context"
	"github.com/Denisz0785/spaceyard/order/internal/model"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.CreateOrderInfo) (*model.CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
	PayOrder(ctx context.Context, orderUUID string, paymentMethod string) error
}
