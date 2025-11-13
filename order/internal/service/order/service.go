package order

import (
	inventoryv1 "github.com/Denisz0785/spaceyard/order/internal/client/grpc"
	"github.com/Denisz0785/spaceyard/order/internal/repo"

	paymentv1 "github.com/Denisz0785/spaceyard/order/internal/client/grpc"
	def "github.com/Denisz0785/spaceyard/order/internal/service"
)

var _ def.OrderService = (*orderService)(nil)

type orderService struct {
	repo repo.OrderRepository

	inventoryClient inventoryv1.InventoryClient
	paymentClient   paymentv1.PaymentClient
}

func NewOrderService(repo repo.OrderRepository, invClient inventoryv1.InventoryClient, payClient paymentv1.PaymentClient) *orderService {
	return &orderService{
		repo: repo,

		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}
