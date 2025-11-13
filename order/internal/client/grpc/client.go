package grpc

import (
	"context"

	"github.com/Denisz0785/spaceyard/order/internal/model"
	"github.com/google/uuid"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (transactionUUID uuid.UUID, err error)
}
