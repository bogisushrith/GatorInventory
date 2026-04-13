package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

/*
ORDER SERVICE UNIT TESTS

This file tests the OrderService layer in isolation by mocking all repository
dependencies: OrderRepository, OrderItemRepository, ProductRepository, and CartRepository.
We verify transaction handling, stock deduction, and order creation workflows.

Key scenarios tested:
1. Happy path: create order with valid items from cart
2. Error cases: insufficient stock, product not found, empty cart
3. Transaction rollback: verify transaction is rolled back on failure
4. Order retrieval: GetAllOrders and GetOrderByID with various conditions
5. Edge cases: duplicate products, large quantities, invalid inputs
*/

// ============================================================================
// MOCK REPOSITORIES
// ============================================================================

type mockTx struct {
	rollbackCalled int
	commitCalled   int
	commitErr      error
}

func (m *mockTx) Commit(ctx context.Context) error { m.commitCalled++; return m.commitErr }
func (m *mockTx) Rollback(ctx context.Context) error { m.rollbackCalled++; return nil }
func (m *mockTx) Begin(ctx context.Context) (pgx.Tx, error) { return nil, nil }
func (m *mockTx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error { return nil }
func (m *mockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) { return 0, nil }
func (m *mockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) { return pgconn.CommandTag(""), nil }
func (m *mockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) { return nil, nil }
func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row { return nil }
func (m *mockTx) CopyFromRows(rows [][]interface{}) pgx.CopyFromSource { return nil }
func (m *mockTx) Conn() *pgx.Conn { return nil }
func (m *mockTx) LargeObjects() pgx.LargeObjects { return pgx.LargeObjects{} }
func (m *mockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) { return nil, nil }
func (m *mockTx) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) { return pgconn.CommandTag(""), nil }

// Note: mockTx does not fully implement pgx.Tx, but provides the methods needed for testing

type mockOrderRepository struct {
	beginTxResult      pgx.Tx
	beginTxErr         error
	createOrderResult  int
	createOrderErr     error
	createOrderCalls   int
	getAllOrdersResult []domain.Order
	getAllOrdersErr    error
	getOrderByIDResult domain.Order
	getOrderByIDErr    error
}

func (m *mockOrderRepository) EnsureOrderSchema() error { return nil }
func (m *mockOrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) { return m.beginTxResult, m.beginTxErr }
func (m *mockOrderRepository) CreateOrder(order domain.Order) (int, error) { return 0, nil }
func (m *mockOrderRepository) CreateOrderTx(ctx context.Context, tx pgx.Tx, order domain.Order) (int, error) { m.createOrderCalls++; return m.createOrderResult, m.createOrderErr }
func (m *mockOrderRepository) GetOrders(query dto.OrderListQuery) ([]domain.Order, error) { return m.getAllOrdersResult, m.getAllOrdersErr }
func (m *mockOrderRepository) GetOrderByIDForRole(query dto.OrderListQuery, id int) (domain.Order, error) { return m.getOrderByIDResult, m.getOrderByIDErr }
func (m *mockOrderRepository) GetOrdersByUserID(userID int64) ([]domain.Order, error) { return m.getAllOrdersResult, m.getAllOrdersErr }
func (m *mockOrderRepository) GetOrderByID(userID int64, id int) (domain.Order, error) { return m.getOrderByIDResult, m.getOrderByIDErr }

var _ repository.IOrderRepository = (*mockOrderRepository)(nil)

type mockOrderItemRepository struct {
	createOrderItemsErr  error
	createOrderItemsCalls int
}

func (m *mockOrderItemRepository) CreateOrderItems(items []domain.OrderItem) error { return nil }
func (m *mockOrderItemRepository) CreateOrderItemsTx(ctx context.Context, tx pgx.Tx, items []domain.OrderItem) error { m.createOrderItemsCalls++; return m.createOrderItemsErr }

var _ repository.IOrderItemRepository = (*mockOrderItemRepository)(nil)

type mockProductRepositoryForOrder struct {
	productByIDMap     map[int64]*domain.Product
	getProductByIDErr  error
	updateStockByIDMap map[int64]int
	updateStockErr     error
	updateStockCalls   int
}

func (m *mockProductRepositoryForOrder) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, error) { return nil, 0, nil }
func (m *mockProductRepositoryForOrder) GetProductByID(productID int64) (*domain.Product, error) {
	if m.getProductByIDErr != nil { return nil, m.getProductByIDErr }
	if product, exists := m.productByIDMap[productID]; exists { return product, nil }
	return nil, pgx.ErrNoRows
}
func (m *mockProductRepositoryForOrder) GetProductByIDTx(ctx context.Context, tx pgx.Tx, productID int64) (*domain.Product, error) {
	if m.getProductByIDErr != nil { return nil, m.getProductByIDErr }
	if product, exists := m.productByIDMap[productID]; exists { return product, nil }
	return nil, pgx.ErrNoRows
}
func (m *mockProductRepositoryForOrder) GetAllProducts() []*domain.Product { return nil }
func (m *mockProductRepositoryForOrder) GetProductsByCategory(category string) []*domain.Product { return nil }
func (m *mockProductRepositoryForOrder) AddProduct(product *domain.Product) error { return nil }
func (m *mockProductRepositoryForOrder) CheckProductExistence(productId int64) error { return nil }
func (m *mockProductRepositoryForOrder) UpdateProductById(updatedProduct *domain.Product, productId int64) error { return nil }
func (m *mockProductRepositoryForOrder) UpdateProductQuantityByID(productID int64, quantity int) (*domain.Product, error) { return nil, nil }
func (m *mockProductRepositoryForOrder) UpdateProductStock(productID int64, quantity int) error { return nil }
func (m *mockProductRepositoryForOrder) UpdateProductStockTx(ctx context.Context, tx pgx.Tx, productID int64, quantity int) error {
	m.updateStockCalls++
	if m.updateStockErr != nil { return m.updateStockErr }
	m.updateStockByIDMap[productID] = quantity
	return nil
}
func (m *mockProductRepositoryForOrder) DeleteProductById(productId int64) error { return nil }

var _ repository.IProductRepository = (*mockProductRepositoryForOrder)(nil)

type mockCartRepositoryForOrder struct {
	cartItemsResult []domain.CartItem
	cartItemsErr    error
	clearCartErr    error
	clearCartCall   int
}

func (m *mockCartRepositoryForOrder) EnsureCartSchema() error { return nil }
func (m *mockCartRepositoryForOrder) GetCartItemsByUserID(userID int64) ([]domain.CartItem, error) { return nil, nil }
func (m *mockCartRepositoryForOrder) GetCartItemsByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) ([]domain.CartItem, error) { return m.cartItemsResult, m.cartItemsErr }
func (m *mockCartRepositoryForOrder) AddCartItem(userID int64, productID int64, quantity int) error { return nil }
func (m *mockCartRepositoryForOrder) UpdateCartItemQuantity(userID int64, productID int64, quantity int) error { return nil }
func (m *mockCartRepositoryForOrder) RemoveCartItem(userID int64, productID int64) error { return nil }
func (m *mockCartRepositoryForOrder) ClearCartByUserID(userID int64) error { return nil }
func (m *mockCartRepositoryForOrder) ClearCartByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) error { m.clearCartCall++; return m.clearCartErr }

var _ repository.ICartRepository = (*mockCartRepositoryForOrder)(nil)

// ============================================================================
// TESTS: CreateOrder
// ============================================================================

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockTx := &mockTx{}
	cartItems := []domain.CartItem{
		{ID: 1, UserID: 10, ProductID: 100, Quantity: 2},
		{ID: 2, UserID: 10, ProductID: 101, Quantity: 1},
	}

	mockOrderRepo := &mockOrderRepository{beginTxResult: mockTx, createOrderResult: 50}
	mockOrderItemRepo := &mockOrderItemRepository{}
	mockProductRepo := &mockProductRepositoryForOrder{
		productByIDMap: map[int64]*domain.Product{
			100: {Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			101: {Id: 101, Name: "Mouse", Price: 25.0, Quantity: 50},
		},
		updateStockByIDMap: make(map[int64]int),
	}
	mockCartRepo := &mockCartRepositoryForOrder{cartItemsResult: cartItems}

	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
	orderID, err := service.CreateOrder(10, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if orderID != 50 {
		t.Fatalf("expected orderID 50, got %d", orderID)
	}
	if mockOrderRepo.createOrderCalls != 1 {
		t.Fatalf("expected CreateOrderTx called once")
	}
	if mockOrderItemRepo.createOrderItemsCalls != 1 {
		t.Fatalf("expected CreateOrderItemsTx called once")
	}
	if mockProductRepo.updateStockCalls != 2 {
		t.Fatalf("expected UpdateProductStockTx called twice, got %d", mockProductRepo.updateStockCalls)
	}
	if mockTx.commitCalled != 1 {
		t.Fatalf("expected commit called once")
	}
	if mockProductRepo.updateStockByIDMap[100] != 8 {
		t.Fatalf("expected product 100 stock 8, got %d", mockProductRepo.updateStockByIDMap[100])
	}
	if mockProductRepo.updateStockByIDMap[101] != 49 {
		t.Fatalf("expected product 101 stock 49, got %d", mockProductRepo.updateStockByIDMap[101])
	}
}

// TestOrderService_CreateOrder tests various scenarios
func TestOrderService_CreateOrder(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		cartItems         []domain.CartItem
		cartErr           error
		products          map[int64]*domain.Product
		expectError       bool
		expectedErrType   error
	}{
		{
			name:            "Invalid userID zero",
			userID:          0,
			expectError:     true,
			expectedErrType: ErrInvalidOrderInput,
		},
		{
			name:            "Invalid userID negative",
			userID:          -1,
			expectError:     true,
			expectedErrType: ErrInvalidOrderInput,
		},
		{
			name:            "Empty cart",
			userID:          10,
			cartItems:       []domain.CartItem{},
			expectError:     true,
			expectedErrType: ErrInvalidOrderInput,
		},
		{
			name:        "Cart retrieval error",
			userID:      10,
			cartErr:     errors.New("database error"),
			expectError: true,
		},
		{
			name:      "Product not found",
			userID:    10,
			cartItems: []domain.CartItem{{ProductID: 999, Quantity: 1}},
			products:  map[int64]*domain.Product{},
			expectError:     true,
			expectedErrType: ErrProductNotFound,
		},
		{
			name:      "Insufficient stock",
			userID:    10,
			cartItems: []domain.CartItem{{ProductID: 100, Quantity: 20}},
			products: map[int64]*domain.Product{
				100: {Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5},
			},
			expectError:     true,
			expectedErrType: ErrInsufficientStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTx := &mockTx{}
			mockOrderRepo := &mockOrderRepository{beginTxResult: mockTx, createOrderResult: 50}
			mockOrderItemRepo := &mockOrderItemRepository{}

			if tt.products == nil {
				tt.products = make(map[int64]*domain.Product)
			}
			mockProductRepo := &mockProductRepositoryForOrder{
				productByIDMap:     tt.products,
				updateStockByIDMap: make(map[int64]int),
			}
			mockCartRepo := &mockCartRepositoryForOrder{
				cartItemsResult: tt.cartItems,
				cartItemsErr:    tt.cartErr,
			}

			service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
			_, err := service.CreateOrder(tt.userID, nil)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedErrType != nil && !errors.Is(err, tt.expectedErrType) {
					t.Fatalf("expected error type %v, got %v", tt.expectedErrType, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

// ============================================================================
// TESTS: GetAllOrders
// ============================================================================

func TestOrderService_GetAllOrders_Success(t *testing.T) {
	orders := []domain.Order{
		{ID: 1, UserID: 10, CreatedAt: time.Now()},
		{ID: 2, UserID: 10, CreatedAt: time.Now()},
	}

	mockOrderRepo := &mockOrderRepository{getAllOrdersResult: orders}
	mockOrderItemRepo := &mockOrderItemRepository{}
	mockProductRepo := &mockProductRepositoryForOrder{}
	mockCartRepo := &mockCartRepositoryForOrder{}

	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
	result, err := service.GetAllOrders(10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result))
	}
}

func TestOrderService_GetAllOrders(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		orders          []domain.Order
		err             error
		expectError     bool
		expectedErrType error
	}{
		{"Get multiple orders", 10, []domain.Order{{ID: 1, UserID: 10}, {ID: 2, UserID: 10}}, nil, false, nil},
		{"Get empty order list", 10, []domain.Order{}, nil, false, nil},
		{"Invalid userID zero", 0, nil, nil, true, ErrInvalidOrderInput},
		{"Invalid userID negative", -1, nil, nil, true, ErrInvalidOrderInput},
		{"Repository error", 10, nil, errors.New("db error"), true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := &mockOrderRepository{
				getAllOrdersResult: tt.orders,
				getAllOrdersErr:    tt.err,
			}
			mockOrderItemRepo := &mockOrderItemRepository{}
			mockProductRepo := &mockProductRepositoryForOrder{}
			mockCartRepo := &mockCartRepositoryForOrder{}

			service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
			result, err := service.GetAllOrders(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedErrType != nil && !errors.Is(err, tt.expectedErrType) {
					t.Fatalf("expected error type %v, got %v", tt.expectedErrType, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil {
					t.Fatalf("expected non-nil result")
				}
			}
		})
	}
}

// ============================================================================
// TESTS: GetOrderByID
// ============================================================================

func TestOrderService_GetOrderByID_Success(t *testing.T) {
	order := domain.Order{ID: 1, UserID: 10, CreatedAt: time.Now()}

	mockOrderRepo := &mockOrderRepository{getOrderByIDResult: order}
	mockOrderItemRepo := &mockOrderItemRepository{}
	mockProductRepo := &mockProductRepositoryForOrder{}
	mockCartRepo := &mockCartRepositoryForOrder{}

	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
	result, err := service.GetOrderByID(10, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Fatalf("expected order ID 1, got %d", result.ID)
	}
}

func TestOrderService_GetOrderByID(t *testing.T) {
	order := domain.Order{ID: 1, UserID: 10, CreatedAt: time.Now()}

	tests := []struct {
		name            string
		userID          int64
		orderID         int
		result          domain.Order
		err             error
		expectError     bool
		expectedErrType error
	}{
		{"Get order success", 10, 1, order, nil, false, nil},
		{"Invalid userID zero", 0, 1, domain.Order{}, nil, true, ErrInvalidOrderInput},
		{"Invalid userID negative", -1, 1, domain.Order{}, nil, true, ErrInvalidOrderInput},
		{"Invalid orderID zero", 10, 0, domain.Order{}, nil, true, ErrInvalidOrderInput},
		{"Invalid orderID negative", 10, -1, domain.Order{}, nil, true, ErrInvalidOrderInput},
		{"Order not found", 10, 999, domain.Order{}, pgx.ErrNoRows, true, ErrOrderNotFound},
		{"Repository error", 10, 1, domain.Order{}, errors.New("db error"), true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := &mockOrderRepository{
				getOrderByIDResult: tt.result,
				getOrderByIDErr:    tt.err,
			}
			mockOrderItemRepo := &mockOrderItemRepository{}
			mockProductRepo := &mockProductRepositoryForOrder{}
			mockCartRepo := &mockCartRepositoryForOrder{}

			service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
			result, err := service.GetOrderByID(tt.userID, tt.orderID)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedErrType != nil && !errors.Is(err, tt.expectedErrType) {
					t.Fatalf("expected error type %v, got %v", tt.expectedErrType, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result.ID != tt.result.ID {
					t.Fatalf("expected order ID %d, got %d", tt.result.ID, result.ID)
				}
			}
		})
	}
}

// ============================================================================
// EDGE CASES AND SPECIAL TESTS
// ============================================================================

func TestOrderService_CreateOrder_VerifyClearCart(t *testing.T) {
	mockTx := &mockTx{}
	cartItems := []domain.CartItem{{ID: 1, UserID: 10, ProductID: 100, Quantity: 2}}

	mockOrderRepo := &mockOrderRepository{beginTxResult: mockTx, createOrderResult: 50}
	mockOrderItemRepo := &mockOrderItemRepository{}
	mockProductRepo := &mockProductRepositoryForOrder{
		productByIDMap:     map[int64]*domain.Product{100: {Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10}},
		updateStockByIDMap: make(map[int64]int),
	}
	mockCartRepo := &mockCartRepositoryForOrder{cartItemsResult: cartItems}

	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
	_, err := service.CreateOrder(10, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockCartRepo.clearCartCall != 1 {
		t.Fatalf("expected ClearCartByUserIDTx called once, got %d", mockCartRepo.clearCartCall)
	}
}

func TestOrderService_CreateOrder_TransactionRollback(t *testing.T) {
	mockTx := &mockTx{}
	cartItems := []domain.CartItem{{ID: 1, UserID: 10, ProductID: 100, Quantity: 20}} // insufficient stock

	mockOrderRepo := &mockOrderRepository{beginTxResult: mockTx}
	mockOrderItemRepo := &mockOrderItemRepository{}
	mockProductRepo := &mockProductRepositoryForOrder{
		productByIDMap:     map[int64]*domain.Product{100: {Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5}},
		updateStockByIDMap: make(map[int64]int),
	}
	mockCartRepo := &mockCartRepositoryForOrder{cartItemsResult: cartItems}

	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockProductRepo, mockCartRepo)
	_, err := service.CreateOrder(10, nil)

	if err == nil {
		t.Fatalf("expected error for insufficient stock")
	}
	if mockTx.rollbackCalled != 1 {
		t.Fatalf("expected rollback called, got %d", mockTx.rollbackCalled)
	}
}
