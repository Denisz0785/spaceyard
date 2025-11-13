package converter

import (
	"github.com/Denisz0785/spaceyard/order/internal/model"
	orderv1 "github.com/Denisz0785/spaceyard/shared/pkg/openapi/order/v1"
)

// OrderToRepoOrder converts a service-level Order model to a repository-level Order model.
// This matches the repository interface which expects *repo/model.Order in Create/Update.
//func OrderToRepoOrder(o *model.Order) *repoModel.Order {
//	if o == nil {
//		return nil
//	}
//	return &repoModel.Order{
//		OrderUUID:       o.OrderUUID,
//		UserUUID:        o.UserUUID,
//		PartUuids:       o.PartUuids,
//		TotalPrice:      o.TotalPrice,
//		TransactionUUID: o.TransactionUUID,
//		PaymentMethod:   o.PaymentMethod,
//		Status:          o.Status,
//	}
//}

func OrderInfoToModel(req *orderv1.CreateOrderRequest) *model.CreateOrderInfo {
	return &model.CreateOrderInfo{
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
}

func ModelToCreateOrderResponse(resp model.CreateOrderResponse) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  resp.OrderUUID,
		TotalPrice: resp.TotalPrice,
	}
}
