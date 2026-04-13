package service

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
)

var (
	ErrInvalidCartInput  = errors.New("invalid cart input")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type ICartService interface {
	GetCart(userID int64) ([]domain.CartItem, error)
	AddToCart(userID int64, productID int64, quantity int) error
	UpdateCartItem(userID int64, productID int64, quantity int) error
	RemoveFromCart(userID int64, productID int64) error
	ClearCart(userID int64) error
}

type CartService struct {
	cartRepository    repository.ICartRepository
	productRepository repository.IProductRepository
}

func NewCartService(cartRepository repository.ICartRepository, productRepository repository.IProductRepository) ICartService {
	return &CartService{cartRepository: cartRepository, productRepository: productRepository}
}

func (service *CartService) GetCart(userID int64) ([]domain.CartItem, error) {
	if userID <= 0 {
		return nil, ErrInvalidCartInput
	}
	return service.cartRepository.GetCartItemsByUserID(userID)
}

func (service *CartService) AddToCart(userID int64, productID int64, quantity int) error {
	if userID <= 0 || productID <= 0 || quantity <= 0 {
		return ErrInvalidCartInput
	}

	product, err := service.productRepository.GetProductByID(productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	if product.Quantity < quantity {
		return ErrInsufficientStock
	}

	return service.cartRepository.AddCartItem(userID, productID, quantity)
}

func (service *CartService) UpdateCartItem(userID int64, productID int64, quantity int) error {
	if userID <= 0 || productID <= 0 {
		return ErrInvalidCartInput
	}

	if quantity > 0 {
		product, err := service.productRepository.GetProductByID(productID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrProductNotFound
			}
			return err
		}
		if product.Quantity < quantity {
			return ErrInsufficientStock
		}
	}

	err := service.cartRepository.UpdateCartItemQuantity(userID, productID, quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCartItemNotFound
		}
		return err
	}
	return nil
}

func (service *CartService) RemoveFromCart(userID int64, productID int64) error {
	if userID <= 0 || productID <= 0 {
		return ErrInvalidCartInput
	}

	err := service.cartRepository.RemoveCartItem(userID, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCartItemNotFound
		}
		return err
	}
	return nil
}

func (service *CartService) ClearCart(userID int64) error {
	if userID <= 0 {
		return ErrInvalidCartInput
	}
	return service.cartRepository.ClearCartByUserID(userID)
}
