package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
	"strings"
)

type IProductRepository interface {
	GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, error)
	GetAllProducts() []*domain.Product
	GetProductsByCategory(category string) []*domain.Product
	AddProduct(product *domain.Product) error
	CheckProductExistence(productId int64) error
	UpdateProductById(updatedProduct *domain.Product, productId int64) error
	DeleteProductById(productId int64) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{dbPool}
}

func (repository *ProductRepository) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, error) {
	ctx := context.Background()

	whereClause, args := buildProductWhereClause(query)

	countStatement := "SELECT COUNT(*) FROM products" + whereClause
	var total int64
	err := repository.dbPool.QueryRow(ctx, countStatement, args...).Scan(&total)
	if err != nil {
		log.Errorf("error while counting products: %v", err)
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	selectStatement := fmt.Sprintf("SELECT id, name, price, quantity, category FROM products%s ORDER BY id DESC LIMIT $%d OFFSET $%d", whereClause, len(args)+1, len(args)+2)
	selectArgs := append(args, query.Limit, offset)

	productRows, err := repository.dbPool.Query(ctx, selectStatement, selectArgs...)
	if err != nil {
		log.Errorf("error while getting products with filters: %v", err)
		return nil, 0, err
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows), total, nil
}

func (repository *ProductRepository) GetAllProducts() []*domain.Product {
	products, _, err := repository.GetProducts(dto.ProductListQuery{Page: 1, Limit: 50})
	if err != nil {
		return nil
	}

	return products
}

func (repository *ProductRepository) GetProductsByCategory(category string) []*domain.Product {
	products, _, err := repository.GetProducts(dto.ProductListQuery{Page: 1, Limit: 50, Category: category})
	if err != nil {
		return nil
	}

	return products
}

func (repository *ProductRepository) AddProduct(product *domain.Product) error {
	ctx := context.Background()

	insertStatement := "INSERT INTO products (name, price, quantity, category) VALUES ($1, $2, $3, $4)"

	addNewProduct, err := repository.dbPool.Exec(ctx, insertStatement, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		log.Errorf("error while adding a new product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product added successfully: %v", addNewProduct))
	return nil
}

func (repository *ProductRepository) CheckProductExistence(productId int64) error {
	ctx := context.Background()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)"
	err := repository.dbPool.QueryRow(ctx, query, productId).Scan(&exists)
	if err != nil {
		log.Errorf("error while checking product existence: %v", err)
		return err
	}

	if !exists {
		return fmt.Errorf("product with id %d does not exist", productId)
	}

	return nil
}

func (repository *ProductRepository) UpdateProductById(updatedProduct *domain.Product, productId int64) error {
	ctx := context.Background()

	updateStatement := "UPDATE products SET name = $1, price = $2, quantity = $3, category = $4 WHERE id = $5"
	_, err := repository.dbPool.Exec(ctx, updateStatement, updatedProduct.Name, updatedProduct.Price, updatedProduct.Quantity, updatedProduct.Category, productId)
	if err != nil {
		log.Errorf("error while updating product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product updated successfully: %v", updatedProduct))
	return nil
}

func (repository *ProductRepository) DeleteProductById(productId int64) error {
	ctx := context.Background()

	deleteExec, err := repository.dbPool.Exec(ctx, "DELETE FROM products WHERE id = $1", productId)
	if err != nil {
		log.Errorf("error while deleting product: %v", err)
		return err
	}

	log.Info("Product deleted successfully")
	log.Info(fmt.Sprintf("%v rows affected", deleteExec.RowsAffected()))

	return nil
}

func extractProductsFromRows(productRows pgx.Rows) []*domain.Product {
	var products []*domain.Product

	for productRows.Next() {
		product := &domain.Product{}
		productRows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity, &product.Category)
		products = append(products, product)
	}

	return products
}

func buildProductWhereClause(query dto.ProductListQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if query.Search != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argIndex))
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	if query.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, query.Category)
		argIndex++
	}

	if query.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *query.MinPrice)
		argIndex++
	}

	if query.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *query.MaxPrice)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}
