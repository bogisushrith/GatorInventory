package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderInput = errors.New("invalid order input")
)

type IOrderService interface {
	CreateOrder(userID int64, orderCreate *dto.OrderCreate) (int, error)
	GetOrders(userID int64, role string, query dto.OrderListQuery) ([]domain.Order, error)
	GetAllOrders(userID int64) ([]domain.Order, error)
	GetOrderByRole(userID int64, role string, id int) (domain.Order, error)
	GetOrderByID(userID int64, id int) (domain.Order, error)
}

type OrderService struct {
	orderRepository     repository.IOrderRepository
	orderItemRepository repository.IOrderItemRepository
	productRepository   repository.IProductRepository
	cartRepository      repository.ICartRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, orderItemRepository repository.IOrderItemRepository, productRepository repository.IProductRepository, cartRepository repository.ICartRepository) IOrderService {
	return &OrderService{orderRepository: orderRepository, orderItemRepository: orderItemRepository, productRepository: productRepository, cartRepository: cartRepository}
}

func (service *OrderService) CreateOrder(userID int64, orderCreate *dto.OrderCreate) (int, error) {
	if userID <= 0 {
		return 0, ErrInvalidOrderInput
	}

	ctx := context.Background()
	tx, err := service.orderRepository.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	cartItems, err := service.cartRepository.GetCartItemsByUserIDTx(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	orderItemsSource := make([]dto.OrderItemCreate, 0)
	if len(cartItems) > 0 {
		for _, cartItem := range cartItems {
			orderItemsSource = append(orderItemsSource, dto.OrderItemCreate{ProductID: int(cartItem.ProductID), Quantity: cartItem.Quantity})
		}
	} else if orderCreate != nil && len(orderCreate.Items) > 0 {
		orderItemsSource = append(orderItemsSource, orderCreate.Items...)
	}

	if len(orderItemsSource) == 0 {
		return 0, ErrInvalidOrderInput
	}
	for _, item := range orderItemsSource {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return 0, ErrInvalidOrderInput
		}
	}

	stockByProductID := map[int]int{}
	for _, item := range orderItemsSource {
		product, getProductErr := service.productRepository.GetProductByIDTx(ctx, tx, int64(item.ProductID))
		if getProductErr != nil {
			if errors.Is(getProductErr, pgx.ErrNoRows) {
				return 0, ErrProductNotFound
			}
			return 0, getProductErr
		}

		updatedStock := product.Quantity - item.Quantity
		if updatedStock < 0 {
			return 0, ErrInsufficientStock
		}

		stockByProductID[item.ProductID] = updatedStock
	}

	orderID, err := service.orderRepository.CreateOrderTx(ctx, tx, domain.Order{UserID: userID, Status: "pending"})
	if err != nil {
		return 0, err
	}

	orderItems := make([]domain.OrderItem, 0, len(orderItemsSource))
	for _, item := range orderItemsSource {
		orderItems = append(orderItems, domain.OrderItem{OrderID: orderID, ProductID: item.ProductID, Quantity: item.Quantity})
	}

	if err = service.orderItemRepository.CreateOrderItemsTx(ctx, tx, orderItems); err != nil {
		return 0, err
	}

	for productID, quantity := range stockByProductID {
		if err = service.productRepository.UpdateProductStockTx(ctx, tx, int64(productID), quantity); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return 0, ErrProductNotFound
			}
			return 0, err
		}
	}

	if err = service.cartRepository.ClearCartByUserIDTx(ctx, tx, userID); err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (service *OrderService) GetAllOrders(userID int64) ([]domain.Order, error) {
	return service.GetOrders(userID, "user", dto.OrderListQuery{UserID: userID, Role: "user"})
}

func (service *OrderService) GetOrderByID(userID int64, id int) (domain.Order, error) {
	return service.GetOrderByRole(userID, "user", id)
}

func (service *OrderService) GetOrders(userID int64, role string, query dto.OrderListQuery) ([]domain.Order, error) {
	if userID <= 0 {
		return nil, ErrInvalidOrderInput
	}

	query.UserID = userID
	query.Role = role

	orders, err := service.orderRepository.GetOrders(query)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (service *OrderService) GetOrderByRole(userID int64, role string, id int) (domain.Order, error) {
	if userID <= 0 || id <= 0 {
		return domain.Order{}, ErrInvalidOrderInput
	}

	order, err := service.orderRepository.GetOrderByIDForRole(dto.OrderListQuery{UserID: userID, Role: role}, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, ErrOrderNotFound
		}
		return domain.Order{}, err
	}

	return order, nil
}
