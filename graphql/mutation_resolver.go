package main

import "context"

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	return nil, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	return nil, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order OrderInput) (*Orders, error) {
	return nil, nil
}
