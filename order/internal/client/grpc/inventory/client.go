package inventory

import (
	"context"
	"fmt"

	"github.com/Denisz0785/spaceyard/order/internal/client/converter"
	clientGrpc "github.com/Denisz0785/spaceyard/order/internal/client/grpc"
	"github.com/Denisz0785/spaceyard/order/internal/model"
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

func (c *inventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]model.Part, error) {
	//partUUIDsStrings := make([]string, len(partUUIDs))
	//for i, u := range partUUIDs {
	//	partUUIDsStrings[i] = u.String()
	//}

	resp, err := c.grpcClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("inventory client: failed to list parts: %w", err)
	}

	return converter.PartsFromProto(resp.GetParts())
}
