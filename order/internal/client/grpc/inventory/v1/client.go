package v1

import (
	clientGrpc "github.com/Denisz0785/spaceyard/order/internal/client/grpc"
	inventoryv1 "github.com/Denisz0785/spaceyard/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

// Проверяем на этапе компиляции, что наша реализация удовлетворяет интерфейсу.
var _ clientGrpc.InventoryClient = (*inventoryClient)(nil)

type inventoryClient struct {
	grpcClient inventoryv1.InventoryServiceClient
}

func New(conn *grpc.ClientConn) *inventoryClient {
	return &inventoryClient{
		grpcClient: inventoryv1.NewInventoryServiceClient(conn),
	}
}
