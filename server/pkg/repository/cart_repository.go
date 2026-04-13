package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

type ICartRepository interface {
	EnsureCartSchema() error
	GetCartItemsByUserID(userID int64) ([]domain.CartItem, error)
	GetCartItemsByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) ([]domain.CartItem, error)
	AddCartItem(userID int64, productID int64, quantity int) error
	UpdateCartItemQuantity(userID int64, productID int64, quantity int) error
	RemoveCartItem(userID int64, productID int64) error
	ClearCartByUserID(userID int64) error
	ClearCartByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) error
}

type CartRepository struct {
	dbPool *pgxpool.Pool
}

func NewCartRepository(dbPool *pgxpool.Pool) ICartRepository {
	return &CartRepository{dbPool: dbPool}
}

func (repository *CartRepository) EnsureCartSchema() error {
	ctx := context.Background()
	statement := `
		CREATE TABLE IF NOT EXISTS cart_items (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			quantity INT NOT NULL CHECK (quantity > 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, product_id)
		);
		CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
		CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);
	`

	_, err := repository.dbPool.Exec(ctx, statement)
	if err != nil {
		log.Errorf("error while ensuring cart schema: %v", err)
		return err
	}

	return nil
}

func (repository *CartRepository) GetCartItemsByUserID(userID int64) ([]domain.CartItem, error) {
	ctx := context.Background()
	return repository.getCartItemsByExecutor(ctx, repository.dbPool, userID)
}

func (repository *CartRepository) GetCartItemsByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) ([]domain.CartItem, error) {
	return repository.getCartItemsByExecutor(ctx, tx, userID)
}

func (repository *CartRepository) AddCartItem(userID int64, productID int64, quantity int) error {
	ctx := context.Background()
	statement := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`

	_, err := repository.dbPool.Exec(ctx, statement, userID, productID, quantity)
	if err != nil {
		log.Errorf("error while adding cart item: %v", err)
		return err
	}
	return nil
}

func (repository *CartRepository) UpdateCartItemQuantity(userID int64, productID int64, quantity int) error {
	ctx := context.Background()
	if quantity <= 0 {
		return repository.RemoveCartItem(userID, productID)
	}

	result, err := repository.dbPool.Exec(ctx, "UPDATE cart_items SET quantity = $1 WHERE user_id = $2 AND product_id = $3", quantity, userID, productID)
	if err != nil {
		log.Errorf("error while updating cart item quantity: %v", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (repository *CartRepository) RemoveCartItem(userID int64, productID int64) error {
	ctx := context.Background()
	result, err := repository.dbPool.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2", userID, productID)
	if err != nil {
		log.Errorf("error while removing cart item: %v", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (repository *CartRepository) ClearCartByUserID(userID int64) error {
	ctx := context.Background()
	_, err := repository.dbPool.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		log.Errorf("error while clearing cart: %v", err)
		return err
	}
	return nil
}

func (repository *CartRepository) ClearCartByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) error {
	_, err := tx.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		log.Errorf("error while clearing cart in transaction: %v", err)
		return err
	}
	return nil
}

func (repository *CartRepository) getCartItemsByExecutor(ctx context.Context, executor DBTX, userID int64) ([]domain.CartItem, error) {
	statement := `
		SELECT c.id, c.user_id, c.product_id, c.quantity, p.name, p.price, p.quantity
		FROM cart_items c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
		ORDER BY c.id ASC`

	rows, err := executor.Query(ctx, statement, userID)
	if err != nil {
		log.Errorf("error while getting cart items: %v", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.CartItem, 0)
	for rows.Next() {
		item := domain.CartItem{}
		if err = rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.ProductName, &item.ProductPrice, &item.ProductStock); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
