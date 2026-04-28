package service

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
	"math"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidQuantity = errors.New("quantity can't be less than zero")
)

type IProductService interface {
	Add(productCreate *dto.ProductCreate) error
	GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, int)
	GetFeaturedProducts(limit int) []*domain.Product
	GetAllProducts() []*domain.Product
	GetAllProductsByCategory(category string) []*domain.Product
	UpdateProductById(updatedProduct *dto.ProductCreate, productId int64) error
	UpdateStockById(productId int64, quantity int) (*domain.Product, error)
	DeleteById(productId int64) error
}

type ProductService struct {
	productRepository repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &ProductService{productRepository}
}

func (service *ProductService) Add(productCreate *dto.ProductCreate) error {
	err := validateProductCreate(productCreate)
	if err != nil {
		return err
	}

	product := productCreateToProduct(productCreate)
	return service.productRepository.AddProduct(product)
}

func (service *ProductService) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, int) {
	products, total, err := service.productRepository.GetProducts(query)
	if err != nil {
		return []*domain.Product{}, 0, 0
	}

	totalPages := 0
	if total > 0 && query.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(query.Limit)))
	}

	return products, total, totalPages
}

func (service *ProductService) GetAllProducts() []*domain.Product {
	return service.productRepository.GetAllProducts()
}

func (service *ProductService) GetFeaturedProducts(limit int) []*domain.Product {
	products, err := service.productRepository.GetFeaturedProducts(limit)
	if err != nil {
		return []*domain.Product{}
	}
	return products
}

func (service *ProductService) GetAllProductsByCategory(category string) []*domain.Product {
	return service.productRepository.GetProductsByCategory(category)
}

func (service *ProductService) UpdateProductById(updatedProduct *dto.ProductCreate, productId int64) error {
	err := validateProductCreate(updatedProduct)
	if err != nil {
		return err
	}

	product := productCreateToProduct(updatedProduct)
	err = service.productRepository.UpdateProductById(product, productId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

func (service *ProductService) UpdateStockById(productId int64, quantity int) (*domain.Product, error) {
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	updatedProduct, err := service.productRepository.UpdateProductQuantityByID(productId, quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return updatedProduct, nil
}

func (service *ProductService) DeleteById(productId int64) error {
	err := service.productRepository.DeleteProductById(productId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

func validateProductCreate(productCreate *dto.ProductCreate) error {
	if productCreate.Name == "" {
		return errors.New("name can't be empty")
	}
	if productCreate.Price < 0 {
		return errors.New("price can't be less than zero")
	}
	if productCreate.Quantity < 0 {
		return ErrInvalidQuantity
	}
	if productCreate.Category == "" {
		return errors.New("category can't be empty")
	}
	return nil
}

func productCreateToProduct(productCreate *dto.ProductCreate) *domain.Product {
	return &domain.Product{
		Name:     productCreate.Name,
		Price:    productCreate.Price,
		Quantity: productCreate.Quantity,
		Category: productCreate.Category,
	}
}
