package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product *Product) error 
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, page, size int) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {

}

func (r *elasticRepository) PutProduct(ctx context.Context, product *Product) error {
	_, err := r.client.Index().
		Index("catalog"). 
		Id(product.ID).
		BodyJson(productDocument{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}).
		Do(ctx)
	if err != nil {
		log.Printf("Failed to index product: %v", err)
		return err
	}
	return nil
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	result, err := r.client.Get().
		Index("catalog").
		Id(id).
		Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	p := productDocument{}
	if err := json.Unmarshal(result.Source, &p); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, page, size int) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(page).Size(size).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Id(id)) // dropped .Type()
	}
	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, doc := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(doc.Source, &p); err != nil {
			return nil, err // BUG 20 fix
		}
		products = append(products, Product{
			ID:          doc.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search(). // BUG 22 fix
					Index("catalog").
					Query(elastic.NewMultiMatchQuery(query, "name", "description")).
					From(int(skip)).Size(int(take)).
					Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			return nil, err // BUG 23 fix
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
