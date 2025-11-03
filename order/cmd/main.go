package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderv1 "github.com/Denisz0785/spaceyard/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Denisz0785/spaceyard/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Denisz0785/spaceyard/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "3080"
	inventoryServiceAddr = "localhost:8080"
	paymentServiceAddr   = "localhost:8081"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// OrderStorage представляет потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderv1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderv1.Order),
	}
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, invClient inventoryv1.InventoryServiceClient, payClient paymentv1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}

// CancelOrder implements cancelOrder operation.
//
// Отменить заказ.
//
// POST /orders/{order_uuid}/cancel
func (s *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	s.storage.mu.Lock() // Lock for the entire read-modify-write operation
	defer s.storage.mu.Unlock()

	orderUUID := params.OrderUUID.String()
	order, ok := s.storage.orders[orderUUID]
	if !ok {
		return &orderv1.CancelOrderNotFound{}, nil
	}

	// Contract: If an order is already PAID, it cannot be cancelled.
	if order.Status == orderv1.OrderStatusPAID {
		return &orderv1.CancelOrderConflict{}, nil
	}

	// If the order is PENDING_PAYMENT, update its status to CANCELLED.
	order.Status = orderv1.OrderStatusCANCELLED
	s.storage.orders[orderUUID] = order

	return &orderv1.CancelOrderNoContent{}, nil
}

// CreateOrder implements createOrder operation.
//
// Создать заказ.
//
// POST /orders
func (s *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	// 1. Call InventoryService to get part details
	partUUIDsStrings := make([]string, len(req.PartUuids))
	for i, u := range req.PartUuids {
		partUUIDsStrings[i] = u.String()
	}

	inventoryResp, err := s.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: partUUIDsStrings,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to communicate with inventory service: %w", err)
	}

	// 2. Validate that all requested parts were found
	if len(inventoryResp.GetParts()) != len(req.PartUuids) {
		return nil, fmt.Errorf("one or more requested parts do not exist")
	}

	// 3. Calculate total price
	var totalPrice float64
	for _, part := range inventoryResp.GetParts() {
		totalPrice += part.GetPrice()
	}

	// 4. Create and save the new order
	order := &orderv1.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: orderv1.OptNilUUID{},
		PaymentMethod:   orderv1.OptPaymentMethod{},
		Status:          orderv1.OrderStatusPENDINGPAYMENT,
	}

	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	s.storage.orders[order.OrderUUID.String()] = order

	// 5. Return the successful response
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

// GetOrder implements getOrder operation.
//
// Получить заказ по UUID.
//
// GET /orders/{order_uuid}
func (s *OrderHandler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()

	order, ok := s.storage.orders[orderUUID]
	if !ok {
		return &orderv1.GetOrderNotFound{}, nil
	}

	return order, nil
}

// PayOrder implements payOrder operation.
// POST /orders/{order_uuid}/pay
func (s *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	s.storage.mu.Lock() // Lock for the entire read-modify-write operation
	defer s.storage.mu.Unlock()

	orderUUID := params.OrderUUID.String()
	order, ok := s.storage.orders[orderUUID]
	if !ok {
		return &orderv1.PayOrderNotFound{}, nil
	}

	// Map OpenAPI PaymentMethod to gRPC PaymentMethod

	paymentMethod, err := toGRPCPaymentMethod(req.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("bad request: %w", err)
	}

	// Call PaymentService to process the payment

	paymentResp, err := s.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call payment service: %w", err)
	}

	// Parse the transaction UUID from the payment service response

	transactionUUID, err := uuid.Parse(paymentResp.GetTransactionUuid())
	if err != nil {
		return nil, fmt.Errorf("payment service returned invalid transaction UUID: %w", err)
	}

	// Update order details with payment information

	order.Status = orderv1.OrderStatusPAID
	order.TransactionUUID = orderv1.NewOptNilUUID(transactionUUID)
	order.PaymentMethod = orderv1.NewOptPaymentMethod(req.PaymentMethod)
	s.storage.orders[orderUUID] = order

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// toGRPCPaymentMethod converts an OpenAPI payment method to a gRPC payment method.

func toGRPCPaymentMethod(method orderv1.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch method {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("unknown payment method: %s", method)
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// --- gRPC Client Setup ---
	invConn, err := grpc.NewClient(inventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := invConn.Close(); cerr != nil {
			log.Printf("Error closing inventory connection: %v", cerr)
		}
	}()

	inventoryClient := inventoryv1.NewInventoryServiceClient(invConn)

	payConn, err := grpc.NewClient(paymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := payConn.Close(); cerr != nil {
			log.Printf("Error closing payment connection: %v", cerr)
		}
	}()
	paymentClient := paymentv1.NewPaymentServiceClient(payConn)
	// --- End gRPC Client Setup ---

	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	srv, err := orderv1.NewServer(orderHandler, orderv1.WithPathPrefix("/api/v1"))
	if err != nil {
		return err
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", srv)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")

	return nil
}
