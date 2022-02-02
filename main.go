package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/DymaSV/shippy-consignment-server/proto/consignment"
	micro "github.com/micro/go-micro/v2"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	mu           sync.RWMutex
	conseignment []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(conseignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	repo.conseignment = append(repo.conseignment, conseignment)
	repo.mu.Unlock()
	return conseignment, nil
}

// Create a new consignment
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.conseignment
}

type consignmentService struct {
	repo repository
}

// Create methode for our service
func (s *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		fmt.Errorf("Cannot add consignment: %v", err)
		return err
	}
	res.Success = true
	res.Consignment = consignment
	return nil
}

func (s *consignmentService) GetConsignment(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}
	service := micro.NewService(
		micro.Name("shippy.consignment.server"),
	)
	service.Init()

	if err := pb.RegisterShippingServiceHandler(service.Server(), &consignmentService{repo}); err != nil {
		log.Panic(err)
	}

	log.Println("Running")
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
