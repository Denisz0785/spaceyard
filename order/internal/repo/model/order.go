package model

import (
	"github.com/Denisz0785/spaceyard/order/internal/model"
	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          model.OrderStatus
}
