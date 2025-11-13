package converter

import (
	"github.com/Denisz0785/spaceyard/order/internal/model"
	repoModel "github.com/Denisz0785/spaceyard/order/internal/repo/model"
)

func OrderToRepoOrder(o *model.Order) *repoModel.Order {
	if o == nil {
		return nil
	}
	return &repoModel.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		Status:          o.Status,
	}
}

func RepoOrderToModel(o *repoModel.Order) *model.Order {
	if o == nil {
		return nil
	}

	return &model.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   o.PaymentMethod,
		Status:          o.Status,
	}
}
