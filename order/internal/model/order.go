package model

import (
	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          OrderStatus
}

type CreateOrderInfo struct {
	UserUUID  uuid.UUID
	PartUuids []uuid.UUID
}

type CreateOrderResponse struct {
	OrderUUID  uuid.UUID
	TotalPrice float64
}
