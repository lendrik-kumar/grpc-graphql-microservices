package main

import "context"

type AccountResolver struct {
	server *Server
}

func (r *AccountResolver) Orders(ctx context.Context, obj *Account) ([]*Orders, error) {
	return obj.Orders, nil
}
