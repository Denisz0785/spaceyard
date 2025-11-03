package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	in "github.com/Denisz0785/spaceyard/shared/pkg/proto/inventory/v1"
)

const (
	port = "localhost:8080"
)

// server is used to implement in.InventoryServiceServer.
type server struct {
	in.UnimplementedInventoryServiceServer
	parts map[string]*in.Part
}

func NewServer() *server {
	return &server{
		UnimplementedInventoryServiceServer: in.UnimplementedInventoryServiceServer{},
		parts:                               make(map[string]*in.Part),
	}
}

// GetPart return part by uuid
func (s server) GetPart(ctx context.Context, req *in.GetPartRequest) (*in.GetPartResponse, error) {
	log.Println("Get request for get part by uuid")

	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "UUID is requred")
	}

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %q not found", req.GetUuid())
	}

	return &in.GetPartResponse{Part: part}, nil
}

type partFilter func(*in.Part) bool

// ListParts returns a list of parts, with optional filtering.
func (s server) ListParts(ctx context.Context, req *in.ListPartsRequest) (*in.ListPartsResponse, error) {
	log.Println("Get request for get list parts by filters")

	filter := req.GetFilter()

	if filter == nil {
		allParts := make([]*in.Part, 0, len(s.parts))
		for _, part := range s.parts {
			allParts = append(allParts, part)
		}
		return &in.ListPartsResponse{Parts: allParts}, nil
	}

	filters := s.buildPartFilters(filter)

	var result []*in.Part
	for _, part := range s.parts {
		matchesAll := true
		for _, f := range filters {
			if !f(part) {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			result = append(result, part)
		}
	}

	return &in.ListPartsResponse{Parts: result}, nil
}

func (s server) buildPartFilters(filter *in.PartsFilter) []partFilter {
	var filters []partFilter

	if len(filter.GetUuids()) > 0 {
		uuidSet := make(map[string]struct{})
		for _, uuid := range filter.GetUuids() {
			uuidSet[uuid] = struct{}{}
		}
		filters = append(filters, func(part *in.Part) bool {
			_, ok := uuidSet[part.GetUuid()]
			return ok
		})
	}

	if len(filter.GetNames()) > 0 {
		nameSet := make(map[string]struct{})
		for _, name := range filter.GetNames() {
			nameSet[name] = struct{}{}
		}
		filters = append(filters, func(part *in.Part) bool {
			_, ok := nameSet[part.GetName()]
			return ok
		})
	}

	if len(filter.GetCategories()) > 0 {
		categorySet := make(map[in.Category]struct{})
		for _, category := range filter.GetCategories() {
			categorySet[category] = struct{}{}
		}
		filters = append(filters, func(part *in.Part) bool {
			_, ok := categorySet[part.GetCategory()]
			return ok
		})
	}

	if len(filter.GetManufacturerCountries()) > 0 {
		countrySet := make(map[string]struct{})
		for _, country := range filter.GetManufacturerCountries() {
			countrySet[country] = struct{}{}
		}
		filters = append(filters, func(part *in.Part) bool {
			if part.GetManufacturer() == nil {
				return false
			}
			_, ok := countrySet[part.GetManufacturer().GetCountry()]
			return ok
		})
	}

	if len(filter.GetTags()) > 0 {
		tagSet := make(map[string]struct{})
		for _, tag := range filter.GetTags() {
			tagSet[tag] = struct{}{}
		}
		filters = append(filters, func(part *in.Part) bool {
			for _, partTag := range part.GetTags() {
				if _, ok := tagSet[partTag]; ok {
					return true
				}
			}
			return false
		})
	}

	return filters
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := NewServer()
	srv.parts["37566f5a-cbb2-49e9-af41-4bc0e49f311a"] = &in.Part{Uuid: "37566f5a-cbb2-49e9-af41-4bc0e49f311a", Name: "star", Price: 450}

	in.RegisterInventoryServiceServer(s, srv)

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
