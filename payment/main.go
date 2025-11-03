package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	po "github.com/Denisz0785/spaceyard/shared/pkg/proto/payment/v1"
)

const (
	port = "localhost:8081"
)

type server struct {
	po.UnimplementedPaymentServiceServer
}

func NewServer() *server {
	return &server{
		UnimplementedPaymentServiceServer: po.UnimplementedPaymentServiceServer{},
	}
}

// PayOrder is doing payment
func (p server) PayOrder(ctx context.Context, req *po.PayOrderRequest) (*po.PayOrderResponse, error) {
	transactionUUID := uuid.New().String()

	// Логируем полученные данные для полноты обработки запроса
	log.Printf(
		"Получен запрос на оплату: OrderUUID=[%s], UserUUID=[%s], PaymentMethod=[%s]",
		req.GetOrderUuid(),
		req.GetUserUuid(),
		req.GetPaymentMethod().String(),
	)

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &po.PayOrderResponse{TransactionUuid: transactionUUID}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := NewServer()

	po.RegisterPaymentServiceServer(s, srv)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	// В конце останавливаем gRPC сервер
	s.GracefulStop()
	log.Println("✅ gRPC server stopped")
}
