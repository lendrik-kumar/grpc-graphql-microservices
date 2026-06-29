package catalog

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PostProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CatalogService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &CatalogService{repository: repository}
}

func (s *CatalogService) PostProduct(ctx context.Context, product *Product) (*Product, error) {
	p := &Product{
		Name:        product.Name,        
		Description: product.Description,
		Price:       product.Price,
		ID:          uuid.New().String(),
	}
	if err := s.repository.PutProduct(ctx, p); err != nil {
		return nil, err 
	}
	return p, nil
}

func (s *CatalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductByID(ctx, id)
}

func (s *CatalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, int(skip), int(take))
}

func (s *CatalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *CatalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
