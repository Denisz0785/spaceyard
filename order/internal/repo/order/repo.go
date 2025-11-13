package order

import (
	def "github.com/Denisz0785/spaceyard/order/internal/repo"
	repoModel "github.com/Denisz0785/spaceyard/order/internal/repo/model"

	"sync"
)

var _ def.OrderRepository = (*storage)(nil)

// orderStorage представляет потокобезопасное хранилище данных о заказах
type storage struct {
	mu     sync.RWMutex
	orders map[string]*repoModel.Order
}

func NewOrderStorage() *storage {
	return &storage{
		orders: make(map[string]*repoModel.Order),
	}
}
