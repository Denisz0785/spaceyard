package v1

import (
	"context"
	"fmt"
	"github.com/Denisz0785/spaceyard/order/internal/client/converter"

	"github.com/Denisz0785/spaceyard/order/internal/model"
	inventoryv1 "github.com/Denisz0785/spaceyard/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]model.Part, error) {
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
