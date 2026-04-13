package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

type IOrderItemRepository interface {
	CreateOrderItems(items []domain.OrderItem) error
	CreateOrderItemsTx(ctx context.Context, tx pgx.Tx, items []domain.OrderItem) error
}

type OrderItemRepository struct {
	dbPool *pgxpool.Pool
}

func NewOrderItemRepository(dbPool *pgxpool.Pool) IOrderItemRepository {
	return &OrderItemRepository{dbPool: dbPool}
}

func (repository *OrderItemRepository) CreateOrderItems(items []domain.OrderItem) error {
	ctx := context.Background()
	for _, item := range items {
		_, err := repository.dbPool.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)", item.OrderID, item.ProductID, item.Quantity)
		if err != nil {
			log.Errorf("error while creating order item: %v", err)
			return err
		}
	}
	return nil
}

func (repository *OrderItemRepository) CreateOrderItemsTx(ctx context.Context, tx pgx.Tx, items []domain.OrderItem) error {
	for _, item := range items {
		_, err := tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)", item.OrderID, item.ProductID, item.Quantity)
		if err != nil {
			log.Errorf("error while creating order item in transaction: %v", err)
			return err
		}
	}
	return nil
}
