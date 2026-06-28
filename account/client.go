package account

import (
	"context"

	"github.com/lendrik-kumar/graphql-grpc-go-microservices/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	r, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Limit: take})
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, len(r.Accounts))
	for i, a := range r.Accounts {
		accounts[i] = &Account{ID: a.Id, Name: a.Name}
	}
	return accounts, nil
}