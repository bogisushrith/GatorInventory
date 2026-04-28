package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"ims-intro/pkg/service/dto"
)

type IAnalyticsRepository interface {
	EnsureSeedData() error
	GetSummary(dateFrom *time.Time, lowStockThreshold int) (dto.AnalyticsSummary, error)
	GetRevenueTrend(dateFrom *time.Time) ([]dto.RevenueTrendPoint, error)
	GetOrderTrend(dateFrom *time.Time) ([]dto.OrderTrendPoint, error)
	GetTopProducts(dateFrom *time.Time, limit int) ([]dto.TopProductMetric, error)
	GetLowStockProducts(threshold int, limit int) ([]dto.LowStockProduct, error)
}

type AnalyticsRepository struct {
	dbPool *pgxpool.Pool
}

func NewAnalyticsRepository(dbPool *pgxpool.Pool) IAnalyticsRepository {
	return &AnalyticsRepository{dbPool: dbPool}
}

func (repository *AnalyticsRepository) EnsureSeedData() error {
	ctx := context.Background()
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var productCount int
	if err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&productCount); err != nil {
		return err
	}

	if productCount == 0 {
		_, err = tx.Exec(ctx, `
			INSERT INTO products (name, price, quantity, category) VALUES
			('Laptop', 1200, 15, 'Electronics'),
			('Mouse', 25, 50, 'Accessories'),
			('Keyboard', 45, 30, 'Accessories'),
			('Monitor', 300, 10, 'Electronics')
		`)
		if err != nil {
			return err
		}
	}

	var orderCount int
	if err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&orderCount); err != nil {
		return err
	}

	if orderCount == 0 {
		var adminID int64
		if err = tx.QueryRow(ctx, "SELECT id FROM users ORDER BY CASE WHEN LOWER(role) = 'admin' THEN 0 ELSE 1 END, id ASC LIMIT 1").Scan(&adminID); err != nil {
			return err
		}

		sampleOrders := []struct {
			productName string
			quantity    int
			createdAt   string
		}{
			{productName: "Laptop", quantity: 2, createdAt: "NOW() - INTERVAL '5 days'"},
			{productName: "Mouse", quantity: 5, createdAt: "NOW() - INTERVAL '3 days'"},
			{productName: "Keyboard", quantity: 3, createdAt: "NOW() - INTERVAL '2 days'"},
			{productName: "Monitor", quantity: 1, createdAt: "NOW() - INTERVAL '1 day'"},
		}

		for _, sampleOrder := range sampleOrders {
			var orderID int
			insertOrderStatement := "INSERT INTO orders (user_id, status, created_at) VALUES ($1, 'completed', " + sampleOrder.createdAt + ") RETURNING id"
			if err = tx.QueryRow(ctx, insertOrderStatement, adminID).Scan(&orderID); err != nil {
				return err
			}

			var productID int64
			if err = tx.QueryRow(ctx, "SELECT id FROM products WHERE name = $1 LIMIT 1", sampleOrder.productName).Scan(&productID); err != nil {
				return err
			}

			if _, err = tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)", orderID, productID, sampleOrder.quantity); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (repository *AnalyticsRepository) GetSummary(dateFrom *time.Time, lowStockThreshold int) (dto.AnalyticsSummary, error) {
	ctx := context.Background()
	query := `
		SELECT
			COALESCE((
				SELECT SUM(oi.quantity * p.price)
				FROM orders o
				JOIN order_items oi ON oi.order_id = o.id
				JOIN products p ON p.id = oi.product_id
				WHERE ($1::timestamp IS NULL OR o.created_at >= $1)
				AND LOWER(COALESCE(o.status, 'pending')) <> 'cancelled'
			), 0) AS total_revenue,
			COALESCE((
				SELECT COUNT(*)
				FROM orders o
				WHERE ($1::timestamp IS NULL OR o.created_at >= $1)
			), 0) AS total_orders,
			COALESCE((SELECT COUNT(*) FROM products), 0) AS total_products,
			COALESCE((
				SELECT COUNT(*)
				FROM products
				WHERE quantity < $2
			), 0) AS low_stock_count
	`

	var summary dto.AnalyticsSummary
	err := repository.dbPool.QueryRow(ctx, query, dateFrom, lowStockThreshold).Scan(
		&summary.TotalRevenue,
		&summary.TotalOrders,
		&summary.TotalProducts,
		&summary.LowStockCount,
	)
	if err != nil {
		return dto.AnalyticsSummary{}, err
	}

	return summary, nil
}

func (repository *AnalyticsRepository) GetRevenueTrend(dateFrom *time.Time) ([]dto.RevenueTrendPoint, error) {
	ctx := context.Background()
	query := `
		SELECT
			DATE(o.created_at) AS metric_date,
			COALESCE(SUM(oi.quantity * p.price), 0) AS value
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON p.id = oi.product_id
		WHERE ($1::timestamp IS NULL OR o.created_at >= $1)
		AND LOWER(COALESCE(o.status, 'pending')) <> 'cancelled'
		GROUP BY DATE(o.created_at)
		ORDER BY DATE(o.created_at) ASC
	`

	rows, err := repository.dbPool.Query(ctx, query, dateFrom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trend := make([]dto.RevenueTrendPoint, 0)
	for rows.Next() {
		var metricDate time.Time
		var value float64
		if err = rows.Scan(&metricDate, &value); err != nil {
			return nil, err
		}
		trend = append(trend, dto.RevenueTrendPoint{Date: metricDate.Format("2006-01-02"), Value: value})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trend, nil
}

func (repository *AnalyticsRepository) GetOrderTrend(dateFrom *time.Time) ([]dto.OrderTrendPoint, error) {
	ctx := context.Background()
	query := `
		SELECT
			DATE(o.created_at) AS metric_date,
			COUNT(*) AS value
		FROM orders o
		WHERE ($1::timestamp IS NULL OR o.created_at >= $1)
		GROUP BY DATE(o.created_at)
		ORDER BY DATE(o.created_at) ASC
	`

	rows, err := repository.dbPool.Query(ctx, query, dateFrom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trend := make([]dto.OrderTrendPoint, 0)
	for rows.Next() {
		var metricDate time.Time
		var value int
		if err = rows.Scan(&metricDate, &value); err != nil {
			return nil, err
		}
		trend = append(trend, dto.OrderTrendPoint{Date: metricDate.Format("2006-01-02"), Value: value})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trend, nil
}

func (repository *AnalyticsRepository) GetTopProducts(dateFrom *time.Time, limit int) ([]dto.TopProductMetric, error) {
	ctx := context.Background()
	query := `
		SELECT
			p.id,
			p.name,
			COALESCE(SUM(oi.quantity), 0) AS total_sold,
			COALESCE(SUM(oi.quantity * p.price), 0) AS revenue
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON p.id = oi.product_id
		WHERE ($1::timestamp IS NULL OR o.created_at >= $1)
		AND LOWER(COALESCE(o.status, 'pending')) <> 'cancelled'
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC, revenue DESC, p.id ASC
		LIMIT $2
	`

	rows, err := repository.dbPool.Query(ctx, query, dateFrom, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]dto.TopProductMetric, 0)
	for rows.Next() {
		var item dto.TopProductMetric
		if err = rows.Scan(&item.ProductID, &item.Name, &item.TotalSold, &item.Revenue); err != nil {
			return nil, err
		}
		products = append(products, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repository *AnalyticsRepository) GetLowStockProducts(threshold int, limit int) ([]dto.LowStockProduct, error) {
	ctx := context.Background()
	query := `
		SELECT id, name, quantity, COALESCE(category, '')
		FROM products
		WHERE quantity < $1
		ORDER BY quantity ASC, id ASC
		LIMIT $2
	`

	rows, err := repository.dbPool.Query(ctx, query, threshold, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]dto.LowStockProduct, 0)
	for rows.Next() {
		var item dto.LowStockProduct
		if err = rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Category); err != nil {
			return nil, err
		}
		products = append(products, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
