package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
)

type IOrderRepository interface {
	EnsureOrderSchema() error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateOrder(order domain.Order) (int, error)
	CreateOrderTx(ctx context.Context, tx pgx.Tx, order domain.Order) (int, error)
	GetOrders(query dto.OrderListQuery) ([]domain.Order, error)
	GetOrderByIDForRole(query dto.OrderListQuery, id int) (domain.Order, error)
	GetOrdersByUserID(userID int64) ([]domain.Order, error)
	GetOrderByID(userID int64, id int) (domain.Order, error)
}

type OrderRepository struct {
	dbPool *pgxpool.Pool
}

func NewOrderRepository(dbPool *pgxpool.Pool) IOrderRepository {
	return &OrderRepository{dbPool: dbPool}
}

func (repository *OrderRepository) EnsureOrderSchema() error {
	ctx := context.Background()
	statement := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		ALTER TABLE orders
		ADD COLUMN IF NOT EXISTS user_id BIGINT;

		ALTER TABLE orders
		ADD COLUMN IF NOT EXISTS status VARCHAR(20);

		UPDATE orders SET status = 'pending' WHERE status IS NULL OR status = '';

		ALTER TABLE orders
		ALTER COLUMN status SET DEFAULT 'pending';

		DELETE FROM order_items
		WHERE order_id IN (SELECT id FROM orders WHERE user_id IS NULL);

		DELETE FROM orders
		WHERE user_id IS NULL;

		ALTER TABLE orders
		ALTER COLUMN user_id SET NOT NULL;

		ALTER TABLE orders
		ALTER COLUMN status SET NOT NULL;

		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id INTEGER NOT NULL REFERENCES products(id),
			quantity INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
		CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	`

	_, err := repository.dbPool.Exec(ctx, statement)
	if err != nil {
		log.Errorf("error while ensuring order schema: %v", err)
		return err
	}

	return nil
}

func (repository *OrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return repository.dbPool.BeginTx(ctx, pgx.TxOptions{})
}

func (repository *OrderRepository) CreateOrder(order domain.Order) (int, error) {
	ctx := context.Background()
	statement := "INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id"

	var orderID int
	status := strings.TrimSpace(order.Status)
	if status == "" {
		status = "pending"
	}
	err := repository.dbPool.QueryRow(ctx, statement, order.UserID, strings.ToLower(status)).Scan(&orderID)
	if err != nil {
		log.Errorf("error while creating order: %v", err)
		return 0, err
	}

	return orderID, nil
}

func (repository *OrderRepository) CreateOrderTx(ctx context.Context, tx pgx.Tx, order domain.Order) (int, error) {
	statement := "INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id"

	var orderID int
	status := strings.TrimSpace(order.Status)
	if status == "" {
		status = "pending"
	}
	err := tx.QueryRow(ctx, statement, order.UserID, strings.ToLower(status)).Scan(&orderID)
	if err != nil {
		log.Errorf("error while creating order in transaction: %v", err)
		return 0, err
	}

	return orderID, nil
}

func (repository *OrderRepository) GetOrdersByUserID(userID int64) ([]domain.Order, error) {
	return repository.GetOrders(dto.OrderListQuery{UserID: userID, Role: "user"})
}

func (repository *OrderRepository) GetOrderByID(userID int64, id int) (domain.Order, error) {
	return repository.GetOrderByIDForRole(dto.OrderListQuery{UserID: userID, Role: "user"}, id)
}

func (repository *OrderRepository) GetOrders(query dto.OrderListQuery) ([]domain.Order, error) {
	ctx := context.Background()
	conditions := []string{"1 = 1"}
	args := make([]interface{}, 0)
	argIndex := 1

	if strings.ToLower(query.Role) != "admin" {
		conditions = append(conditions, fmt.Sprintf("o.user_id = $%d", argIndex))
		args = append(args, query.UserID)
		argIndex++
	}

	if query.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(CAST(o.id AS TEXT) ILIKE $%d OR LOWER(COALESCE(u.username, '')) ILIKE LOWER($%d))", argIndex, argIndex))
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	if query.Status != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(o.status, 'pending')) = LOWER($%d)", argIndex))
		args = append(args, query.Status)
		argIndex++
	}

	if query.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("o.created_at >= $%d", argIndex))
		args = append(args, query.DateFrom.UTC())
		argIndex++
	}

	if query.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("o.created_at < $%d", argIndex))
		args = append(args, query.DateTo.UTC())
		argIndex++
	}

	statement := fmt.Sprintf(`
		SELECT o.id, o.user_id, COALESCE(u.username, ''), COALESCE(o.status, 'pending'), o.created_at,
		       oi.id, oi.order_id, oi.product_id, COALESCE(p.name, ''), COALESCE(p.price, 0), oi.quantity
		FROM orders o
		LEFT JOIN users u ON u.id = o.user_id
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE %s
		ORDER BY o.created_at DESC, o.id DESC, oi.id ASC`, strings.Join(conditions, " AND "))

	rows, err := repository.dbPool.Query(ctx, statement, args...)
	if err != nil {
		log.Errorf("error while fetching orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	return assembleOrdersFromRows(rows)
}

func (repository *OrderRepository) GetOrderByIDForRole(query dto.OrderListQuery, id int) (domain.Order, error) {
	ctx := context.Background()
	conditions := []string{"o.id = $1"}
	args := []interface{}{id}

	if strings.ToLower(query.Role) != "admin" {
		conditions = append(conditions, fmt.Sprintf("o.user_id = $%d", len(args)+1))
		args = append(args, query.UserID)
	}

	statement := fmt.Sprintf(`
		SELECT o.id, o.user_id, COALESCE(u.username, ''), COALESCE(o.status, 'pending'), o.created_at,
		       oi.id, oi.order_id, oi.product_id, COALESCE(p.name, ''), COALESCE(p.price, 0), oi.quantity
		FROM orders o
		LEFT JOIN users u ON u.id = o.user_id
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE %s
		ORDER BY oi.id ASC`, strings.Join(conditions, " AND "))

	rows, err := repository.dbPool.Query(ctx, statement, args...)
	if err != nil {
		log.Errorf("error while fetching order %d: %v", id, err)
		return domain.Order{}, err
	}
	defer rows.Close()

	orders, err := assembleOrdersFromRows(rows)
	if err != nil {
		return domain.Order{}, err
	}
	if len(orders) == 0 {
		return domain.Order{}, pgx.ErrNoRows
	}
	return orders[0], nil
}

func assembleOrdersFromRows(rows pgx.Rows) ([]domain.Order, error) {
	ordersMap := map[int]*domain.Order{}
	orderIDs := make([]int, 0)

	for rows.Next() {
		var orderID int
		var orderUserID int64
		var userName string
		var status string
		var createdAt time.Time
		var itemID, itemOrderID, productID, quantity *int
		var productName *string
		var productPrice *float32

		err := rows.Scan(&orderID, &orderUserID, &userName, &status, &createdAt, &itemID, &itemOrderID, &productID, &productName, &productPrice, &quantity)
		if err != nil {
			return nil, err
		}

		if _, exists := ordersMap[orderID]; !exists {
			ordersMap[orderID] = &domain.Order{ID: orderID, UserID: orderUserID, UserName: userName, Status: status, CreatedAt: createdAt, Items: []domain.OrderItem{}}
			orderIDs = append(orderIDs, orderID)
		}

		if itemID != nil && itemOrderID != nil && productID != nil && quantity != nil {
			item := domain.OrderItem{
				ID:          *itemID,
				OrderID:     *itemOrderID,
				ProductID:   *productID,
				ProductName: "",
				Quantity:    *quantity,
			}
			if productName != nil {
				item.ProductName = *productName
			}
			if productPrice != nil {
				item.ProductPrice = *productPrice
			}
			ordersMap[orderID].Items = append(ordersMap[orderID].Items, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(orderIDs, func(i, j int) bool { return orderIDs[i] > orderIDs[j] })
	orders := make([]domain.Order, 0, len(orderIDs))
	for _, id := range orderIDs {
		orders = append(orders, *ordersMap[id])
	}

	return orders, nil
}
