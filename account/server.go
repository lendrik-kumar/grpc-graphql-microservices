package account

import (
	"context"
	"fmt"
	"net"

	"github.com/lendrik-kumar/graphql-grpc-go-microservices/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenAndServeGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{service: service})
	reflection.Register(server)

	return server.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountRespone, error) {
	account, err := s.service.PostAccount(ctx, &Account{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountRespone{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReposne, error) {
	account, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountReposne{Account: &pb.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsRespone, error) {
	accounts, err := s.service.GetAccounts(ctx, req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}

	pbAccounts := make([]*pb.Account, len(accounts))
	for i, account := range accounts {
		pbAccounts[i] = &pb.Account{Id: account.ID, Name: account.Name}
	}
	return &pb.GetAccountsRespone{Accounts: pbAccounts}, nil
}
