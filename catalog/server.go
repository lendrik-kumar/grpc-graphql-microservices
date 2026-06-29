package catalog

import (
	"context"
	"net"

	"github.com/lendrik-kumar/graphql-grpc-go-microservices/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenAndServeGRPC(service Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterCatalogServiceServer(server, &grpcServer{
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var response []Product
	var err error
	if req.Query != "" {
		response, err = s.service.SearchProducts(ctx, req.Query, req.Skip, req.Take)
	} else if len(req.Ids) > 0 {
		response, err = s.service.GetProductsByIDs(ctx, req.Ids)
	} else {
		response, err = s.service.GetProducts(ctx, req.Skip, req.Take)
	}
	if err != nil {
		return nil, err
	}
	products := []*pb.Product{}
	for _, p := range response {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
