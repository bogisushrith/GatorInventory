package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"ims-intro/pkg/service/dto"
)

type IUserAnalyticsRepository interface {
	GetSummary(userID int64) (dto.UserAnalyticsSummary, error)
	GetRecentOrders(userID int64, limit int) ([]dto.UserRecentOrder, error)
	GetTopProducts(userID int64, limit int) ([]dto.UserTopProduct, error)
	GetSpendingTrend(userID int64) ([]dto.RevenueTrendPoint, error)
	GetRecommendations(userID int64, limit int) ([]dto.UserRecommendation, error)
}

type UserAnalyticsRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserAnalyticsRepository(dbPool *pgxpool.Pool) IUserAnalyticsRepository {
	return &UserAnalyticsRepository{dbPool: dbPool}
}

func (repository *UserAnalyticsRepository) GetSummary(userID int64) (dto.UserAnalyticsSummary, error) {
	ctx := context.Background()
	query := `
		SELECT
			COUNT(DISTINCT o.id) AS total_orders,
			COALESCE(SUM(oi.quantity * p.price), 0) AS total_spent,
			COALESCE(SUM(CASE WHEN LOWER(COALESCE(o.status, 'pending')) = 'pending' THEN 1 ELSE 0 END), 0) AS pending_orders
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.user_id = $1
	`

	var summary dto.UserAnalyticsSummary
	err := repository.dbPool.QueryRow(ctx, query, userID).Scan(&summary.TotalOrders, &summary.TotalSpent, &summary.PendingOrders)
	if err != nil {
		return dto.UserAnalyticsSummary{}, err
	}

	return summary, nil
}

func (repository *UserAnalyticsRepository) GetRecentOrders(userID int64, limit int) ([]dto.UserRecentOrder, error) {
	ctx := context.Background()
	query := `
		SELECT
			o.id,
			COALESCE(STRING_AGG(DISTINCT p.name, ', ' ORDER BY p.name), 'No items') AS product_names,
			LOWER(COALESCE(o.status, 'pending')) AS status,
			o.created_at,
			COALESCE(SUM(oi.quantity * p.price), 0) AS total_price
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE o.user_id = $1
		GROUP BY o.id, o.status, o.created_at
		ORDER BY o.created_at DESC, o.id DESC
		LIMIT $2
	`

	rows, err := repository.dbPool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dto.UserRecentOrder, 0)
	for rows.Next() {
		var item dto.UserRecentOrder
		if err = rows.Scan(&item.OrderID, &item.ProductNames, &item.Status, &item.CreatedAt, &item.TotalPrice); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (repository *UserAnalyticsRepository) GetTopProducts(userID int64, limit int) ([]dto.UserTopProduct, error) {
	ctx := context.Background()
	query := `
		SELECT
			p.id,
			p.name,
			COUNT(DISTINCT oi.order_id) AS purchase_count
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON p.id = oi.product_id
		WHERE o.user_id = $1
		GROUP BY p.id, p.name
		ORDER BY purchase_count DESC, p.name ASC
		LIMIT $2
	`

	rows, err := repository.dbPool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dto.UserTopProduct, 0)
	for rows.Next() {
		var item dto.UserTopProduct
		if err = rows.Scan(&item.ProductID, &item.Name, &item.PurchaseCount); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func (repository *UserAnalyticsRepository) GetSpendingTrend(userID int64) ([]dto.RevenueTrendPoint, error) {
	ctx := context.Background()
	query := `
		SELECT
			DATE(o.created_at) AS metric_date,
			COALESCE(SUM(oi.quantity * p.price), 0) AS total_spent
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON p.id = oi.product_id
		WHERE o.user_id = $1
		GROUP BY DATE(o.created_at)
		ORDER BY DATE(o.created_at) ASC
	`

	rows, err := repository.dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dto.RevenueTrendPoint, 0)
	for rows.Next() {
		var metricDate time.Time
		var value float64
		if err = rows.Scan(&metricDate, &value); err != nil {
			return nil, err
		}
		result = append(result, dto.RevenueTrendPoint{
			Date:  metricDate.Format("2006-01-02"),
			Value: value,
		})
	}

	return result, rows.Err()
}

func (repository *UserAnalyticsRepository) GetRecommendations(userID int64, limit int) ([]dto.UserRecommendation, error) {
	ctx := context.Background()
	query := `
		WITH user_categories AS (
			SELECT p.category, SUM(oi.quantity) AS total_quantity
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			JOIN products p ON p.id = oi.product_id
			WHERE o.user_id = $1
			GROUP BY p.category
			ORDER BY total_quantity DESC
			LIMIT 3
		),
		user_products AS (
			SELECT DISTINCT oi.product_id
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE o.user_id = $1
		),
		similar_products AS (
			SELECT
				p.id,
				p.name,
				p.price,
				p.category,
				'Recommended from your favorite categories' AS reason,
				COALESCE(SUM(oi.quantity), 0) AS popularity
			FROM products p
			LEFT JOIN order_items oi ON oi.product_id = p.id
			WHERE p.category IN (SELECT category FROM user_categories)
			AND p.id NOT IN (SELECT product_id FROM user_products)
			GROUP BY p.id, p.name, p.price, p.category
		),
		popular_products AS (
			SELECT
				p.id,
				p.name,
				p.price,
				p.category,
				'Popular with shoppers right now' AS reason,
				COALESCE(SUM(oi.quantity), 0) AS popularity
			FROM products p
			LEFT JOIN order_items oi ON oi.product_id = p.id
			GROUP BY p.id, p.name, p.price, p.category
		)
		SELECT id, name, price, category, reason
		FROM (
			SELECT * FROM similar_products
			UNION ALL
			SELECT * FROM popular_products
		) combined
		GROUP BY id, name, price, category, reason, popularity
		ORDER BY CASE WHEN reason = 'Recommended from your favorite categories' THEN 0 ELSE 1 END, popularity DESC, name ASC
		LIMIT $2
	`

	rows, err := repository.dbPool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dto.UserRecommendation, 0)
	for rows.Next() {
		var item dto.UserRecommendation
		if err = rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Category, &item.Reason); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}
