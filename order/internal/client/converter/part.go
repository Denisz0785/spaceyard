package converter

import (
	"github.com/Denisz0785/spaceyard/order/internal/model"
	inventoryv1 "github.com/Denisz0785/spaceyard/shared/pkg/proto/inventory/v1"
)

// PartFromProto преобразует protobuf-модель Part в доменную модель Part.
func PartFromProto(part *inventoryv1.Part) (model.Part, error) {
	// Конвертируем вложенные структуры. Protobuf-версии могут быть nil.
	var dims model.Dimensions
	if pDims := part.GetDimensions(); pDims != nil {
		dims = DimensionsFromProto(pDims)
	}

	var manufacturer model.Manufacturer
	if pMan := part.GetManufacturer(); pMan != nil {
		manufacturer = ManufacturerFromProto(pMan)
	}

	// Конвертируем карту метаданных.
	metadata := make(map[string]model.Value)
	if pMeta := part.GetMetadata(); pMeta != nil {
		for k, v := range pMeta {
			metadata[k] = ValueFromProto(v)
		}
	}

	return model.Part{
		UUID:          part.GetUuid(), // Прямое присваивание, так как оба - строки
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      model.Category(part.GetCategory()), // Простое приведение типов для enum
		Dimensions:    dims,
		Manufacturer:  manufacturer,
		Tags:          part.GetTags(), // Слайсы строк идентичны
		Metadata:      metadata,
		CreatedAt:     part.GetCreatedAt().AsTime(),
		UpdatedAt:     part.GetUpdatedAt().AsTime(),
	}, nil
}

// PartsFromProto преобразует срез protobuf-моделей Part в срез доменных моделей.
func PartsFromProto(parts []*inventoryv1.Part) ([]model.Part, error) {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		p, err := PartFromProto(part)
		if err != nil {
			// В реальном приложении здесь может быть более сложная логика обработки ошибок
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// DimensionsFromProto преобразует protobuf-модель Dimensions в доменную.
func DimensionsFromProto(dims *inventoryv1.Dimensions) model.Dimensions {
	if dims == nil {
		return model.Dimensions{}
	}
	return model.Dimensions{
		Length: dims.GetLength(),
		Width:  dims.GetWidth(),
		Height: dims.GetHeight(),
		Weight: dims.GetWeight(),
	}
}

// ManufacturerFromProto преобразует protobuf-модель Manufacturer в доменную.
func ManufacturerFromProto(m *inventoryv1.Manufacturer) model.Manufacturer {
	if m == nil {
		return model.Manufacturer{}
	}
	return model.Manufacturer{
		Name:    m.GetName(),
		Country: m.GetCountry(),
		Website: m.GetWebsite(),
	}
}

// ValueFromProto преобразует protobuf-модель Value (oneof) в доменную.
func ValueFromProto(v *inventoryv1.Value) model.Value {
	if v == nil {
		return model.Value{}
	}
	switch val := v.Value.(type) {
	case *inventoryv1.Value_StringValue:
		return model.Value{String: val.StringValue}
	case *inventoryv1.Value_Int64Value:
		return model.Value{Int64: val.Int64Value}
	case *inventoryv1.Value_DoubleValue:
		return model.Value{Double: val.DoubleValue}
	case *inventoryv1.Value_BoolValue:
		return model.Value{Bool: val.BoolValue}
	default:
		return model.Value{}
	}
}
