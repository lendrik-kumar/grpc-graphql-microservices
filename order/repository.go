package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID, o.CreatedAt, o.AccountID, o.TotalPrice,
	)
	if err != nil {
		return err
	}

	// fixed: error no longer ignored
	stm, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		return err
	}
	defer stm.Close()

	for _, p := range o.Products {
		// fixed: use = not := so outer err is set (triggers rollback on failure)
		_, err = stm.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id, o.created_at, o.account_id, o.total_price::money::numeric::float8, op.product_id, op.quantity
		 FROM orders o
		 JOIN order_products op ON (o.id = op.order_id)
		 WHERE o.account_id = $1
		 ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	lastOrder := &Order{}
	products := []OrderedProduct{}

	for rows.Next() {
		var p OrderedProduct
		var currentID string

		// fixed: single scan per row into a temp currentID + p, not double-scanned
		if err := rows.Scan(
			&currentID,
			&lastOrder.CreatedAt,
			&lastOrder.AccountID,
			&lastOrder.TotalPrice,
			&p.ID,
			&p.Quantity,
		); err != nil {
			return nil, err
		}

		// flush previous order group when ID changes
		if lastOrder.ID != "" && lastOrder.ID != currentID {
			lastOrder.Products = products // fixed: assign products before flushing
			orders = append(orders, *lastOrder)
			products = []OrderedProduct{}
		}

		lastOrder.ID = currentID
		products = append(products, p)
	}

	// flush the final order group
	if lastOrder.ID != "" {
		lastOrder.Products = products
		orders = append(orders, *lastOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
