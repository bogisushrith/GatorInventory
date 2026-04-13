package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

/*
CART SERVICE UNIT TESTS

This file tests the CartService layer in isolation by mocking the CartRepository
and ProductRepository dependencies. We verify:
1. Input validation (invalid IDs, quantities)
2. Business logic (stock checks, cart operations)
3. Error handling (product not found, insufficient stock)
4. Happy paths (add, update, remove, clear cart items)

The mock repositories simulate database success/failure scenarios without
hitting the actual PostgreSQL database, making tests fast and deterministic.
*/

// ============================================================================
// MOCK REPOSITORIES (for CartService testing)
// ============================================================================

type mockCartRepository struct {
	// For GetCartItemsByUserID
	getCartItemsResult []domain.CartItem
	getCartItemsErr    error
	getCartUserIDCalls int

	// For AddCartItem
	addCartItemErr error
	addCartItemCalls struct {
		calls int
		args  []struct {
			userID    int64
			productID int64
			quantity  int
		}
	}

	// For UpdateCartItemQuantity
	updateCartItemErr error
	updateCartItemCalls struct {
		calls int
		args  []struct {
			userID    int64
			productID int64
			quantity  int
		}
	}

	// For RemoveCartItem
	removeCartItemErr error
	removeCartItemCalls struct {
		calls int
		args  []struct {
			userID    int64
			productID int64
		}
	}

	// For ClearCartByUserID
	clearCartErr error
	clearCartCalls int
}

func (m *mockCartRepository) EnsureCartSchema() error {
	return nil
}

func (m *mockCartRepository) GetCartItemsByUserID(userID int64) ([]domain.CartItem, error) {
	m.getCartUserIDCalls++
	return m.getCartItemsResult, m.getCartItemsErr
}

func (m *mockCartRepository) GetCartItemsByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) ([]domain.CartItem, error) {
	return m.getCartItemsResult, m.getCartItemsErr
}

func (m *mockCartRepository) AddCartItem(userID int64, productID int64, quantity int) error {
	m.addCartItemCalls.calls++
	m.addCartItemCalls.args = append(m.addCartItemCalls.args, struct {
		userID    int64
		productID int64
		quantity  int
	}{userID, productID, quantity})
	return m.addCartItemErr
}

func (m *mockCartRepository) UpdateCartItemQuantity(userID int64, productID int64, quantity int) error {
	m.updateCartItemCalls.calls++
	m.updateCartItemCalls.args = append(m.updateCartItemCalls.args, struct {
		userID    int64
		productID int64
		quantity  int
	}{userID, productID, quantity})
	return m.updateCartItemErr
}

func (m *mockCartRepository) RemoveCartItem(userID int64, productID int64) error {
	m.removeCartItemCalls.calls++
	m.removeCartItemCalls.args = append(m.removeCartItemCalls.args, struct {
		userID    int64
		productID int64
	}{userID, productID})
	return m.removeCartItemErr
}

func (m *mockCartRepository) ClearCartByUserID(userID int64) error {
	m.clearCartCalls++
	return m.clearCartErr
}

func (m *mockCartRepository) ClearCartByUserIDTx(ctx context.Context, tx pgx.Tx, userID int64) error {
	return m.clearCartErr
}

// Verify mockCartRepository implements ICartRepository
var _ repository.ICartRepository = (*mockCartRepository)(nil)

// ============================================================================
// MOCK PRODUCT REPOSITORY (for CartService testing)
// ============================================================================

type mockProductRepositoryForCart struct {
	getProductByIDResult *domain.Product
	getProductByIDErr    error
	getProductByIDCalls  int
}

func (m *mockProductRepositoryForCart) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockProductRepositoryForCart) GetProductByID(productID int64) (*domain.Product, error) {
	m.getProductByIDCalls++
	return m.getProductByIDResult, m.getProductByIDErr
}

func (m *mockProductRepositoryForCart) GetProductByIDTx(ctx context.Context, tx pgx.Tx, productID int64) (*domain.Product, error) {
	return m.getProductByIDResult, m.getProductByIDErr
}

func (m *mockProductRepositoryForCart) GetAllProducts() []*domain.Product {
	return nil
}

func (m *mockProductRepositoryForCart) GetProductsByCategory(category string) []*domain.Product {
	return nil
}

func (m *mockProductRepositoryForCart) AddProduct(product *domain.Product) error {
	return nil
}

func (m *mockProductRepositoryForCart) CheckProductExistence(productId int64) error {
	return nil
}

func (m *mockProductRepositoryForCart) UpdateProductById(updatedProduct *domain.Product, productId int64) error {
	return nil
}

func (m *mockProductRepositoryForCart) UpdateProductQuantityByID(productID int64, quantity int) (*domain.Product, error) {
	return nil, nil
}

func (m *mockProductRepositoryForCart) UpdateProductStock(productID int64, quantity int) error {
	return nil
}

func (m *mockProductRepositoryForCart) UpdateProductStockTx(ctx context.Context, tx pgx.Tx, productID int64, quantity int) error {
	return nil
}

func (m *mockProductRepositoryForCart) DeleteProductById(productId int64) error {
	return nil
}

// Verify mockProductRepositoryForCart implements IProductRepository
var _ repository.IProductRepository = (*mockProductRepositoryForCart)(nil)

// ============================================================================
// TESTS: GetCart
// ============================================================================

// TestCartService_GetCart_Success verifies cart retrieval with valid items.
// Why important: GetCart is the primary read operation for user's cart state.
// Validates: returns correct cart items from repository.
// Setup: mock repo returns list of cart items for valid user.
// Execution: call GetCart with valid userID.
// Assertion: expect no error and correct cart items returned.
func TestCartService_GetCart_Success(t *testing.T) {
	mockCartRepo := &mockCartRepository{
		getCartItemsResult: []domain.CartItem{
			{ID: 1, UserID: 10, ProductID: 100, Quantity: 2, ProductName: "Laptop", ProductPrice: 1500.0},
			{ID: 2, UserID: 10, ProductID: 101, Quantity: 1, ProductName: "Mouse", ProductPrice: 25.0},
		},
	}
	mockProductRepo := &mockProductRepositoryForCart{}

	service := NewCartService(mockCartRepo, mockProductRepo)
	cartItems, err := service.GetCart(10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cartItems) != 2 {
		t.Fatalf("expected 2 items, got %d", len(cartItems))
	}
	if cartItems[0].ProductName != "Laptop" {
		t.Fatalf("expected product name Laptop, got %s", cartItems[0].ProductName)
	}
	if mockCartRepo.getCartUserIDCalls != 1 {
		t.Fatalf("expected GetCartItemsByUserID called once, got %d times", mockCartRepo.getCartUserIDCalls)
	}
}

// TestCartService_GetCart_EmptyCart verifies empty cart is returned correctly.
// Why important: empty cart is valid state and must be handled gracefully.
// Validates: returns empty slice (not nil) for user with no items.
// Setup: mock repo returns empty slice for valid user.
// Execution: call GetCart with valid userID.
// Assertion: expect no error and empty slice.
func TestCartService_GetCart_EmptyCart(t *testing.T) {
	mockCartRepo := &mockCartRepository{
		getCartItemsResult: []domain.CartItem{},
	}
	mockProductRepo := &mockProductRepositoryForCart{}

	service := NewCartService(mockCartRepo, mockProductRepo)
	cartItems, err := service.GetCart(10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cartItems) != 0 {
		t.Fatalf("expected 0 items, got %d", len(cartItems))
	}
}

// TestCartService_GetCart_InvalidUserID validates input validation.
// Why important: prevents invalid database queries.
// Validates: service rejects invalid user IDs without calling repository.
// Setup: no repo calls expected.
// Execution: call GetCart with invalid userID.
// Assertion: expect error and no repository interaction.
func TestCartService_GetCart_InvalidUserID(t *testing.T) {
	tests := []struct {
		name   string
		userID int64
	}{
		{"negative ID", -1},
		{"zero ID", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := &mockCartRepository{}
			mockProductRepo := &mockProductRepositoryForCart{}

			service := NewCartService(mockCartRepo, mockProductRepo)
			_, err := service.GetCart(tt.userID)

			if err == nil {
				t.Fatalf("expected error for userID %d, got nil", tt.userID)
			}
			if !errors.Is(err, ErrInvalidCartInput) {
				t.Fatalf("expected ErrInvalidCartInput, got %v", err)
			}
			if mockCartRepo.getCartUserIDCalls != 0 {
				t.Fatalf("expected no repository calls, got %d", mockCartRepo.getCartUserIDCalls)
			}
		})
	}
}

// ============================================================================
// TESTS: AddToCart
// ============================================================================

// TestCartService_AddToCart table-driven tests for various scenarios.
func TestCartService_AddToCart(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		productID         int64
		quantity          int
		product           *domain.Product
		productErr        error
		addErr            error
		expectError       bool
		expectedErrType   error
		shouldCallAddCart bool
	}{
		// Success case
		{
			name:              "Add new product to cart success",
			userID:            10,
			productID:         100,
			quantity:          2,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5},
			expectError:       false,
			shouldCallAddCart: true,
		},
		// Duplicate product (incrementing quantity)
		{
			name:              "Add same product again increments quantity",
			userID:            10,
			productID:         100,
			quantity:          3,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			expectError:       false,
			shouldCallAddCart: true,
		},
		// Invalid inputs
		{
			name:              "Invalid userID negative",
			userID:            -1,
			productID:         100,
			quantity:          1,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallAddCart: false,
		},
		{
			name:              "Invalid userID zero",
			userID:            0,
			productID:         100,
			quantity:          1,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallAddCart: false,
		},
		{
			name:              "Invalid productID negative",
			userID:            10,
			productID:         -1,
			quantity:          1,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallAddCart: false,
		},
		{
			name:              "Invalid quantity zero",
			userID:            10,
			productID:         100,
			quantity:          0,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallAddCart: false,
		},
		{
			name:              "Invalid quantity negative",
			userID:            10,
			productID:         100,
			quantity:          -5,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallAddCart: false,
		},
		// Product not found
		{
			name:              "Product not found",
			userID:            10,
			productID:         999,
			quantity:          1,
			productErr:        pgx.ErrNoRows,
			expectError:       true,
			expectedErrType:   ErrProductNotFound,
			shouldCallAddCart: false,
		},
		// Generic product repository error
		{
			name:              "Product repository error",
			userID:            10,
			productID:         100,
			quantity:          1,
			productErr:        errors.New("database error"),
			expectError:       true,
			shouldCallAddCart: false,
		},
		// Insufficient stock
		{
			name:              "Insufficient stock",
			userID:            10,
			productID:         100,
			quantity:          10,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 3},
			expectError:       true,
			expectedErrType:   ErrInsufficientStock,
			shouldCallAddCart: false,
		},
		// AddCartItem repository error
		{
			name:              "AddCartItem repository error",
			userID:            10,
			productID:         100,
			quantity:          2,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5},
			addErr:            errors.New("database error"),
			expectError:       true,
			shouldCallAddCart: true,
		},
		// Edge case: exactly matching stock
		{
			name:              "Quantity exactly matches available stock",
			userID:            10,
			productID:         100,
			quantity:          5,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5},
			expectError:       false,
			shouldCallAddCart: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := &mockCartRepository{addCartItemErr: tt.addErr}
			mockProductRepo := &mockProductRepositoryForCart{
				getProductByIDResult: tt.product,
				getProductByIDErr:    tt.productErr,
			}

			service := NewCartService(mockCartRepo, mockProductRepo)
			err := service.AddToCart(tt.userID, tt.productID, tt.quantity)

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

			if tt.shouldCallAddCart {
				if mockCartRepo.addCartItemCalls.calls != 1 {
					t.Fatalf("expected AddCartItem called once, got %d", mockCartRepo.addCartItemCalls.calls)
				}
			} else {
				if mockCartRepo.addCartItemCalls.calls != 0 {
					t.Fatalf("expected AddCartItem not called, got %d calls", mockCartRepo.addCartItemCalls.calls)
				}
			}
		})
	}
}

// ============================================================================
// TESTS: UpdateCartItem
// ============================================================================

// TestCartService_UpdateCartItem table-driven tests for quantity updates.
func TestCartService_UpdateCartItem(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		productID         int64
		quantity          int
		product           *domain.Product
		productErr        error
		updateErr         error
		expectError       bool
		expectedErrType   error
		shouldCallProduct bool
		shouldCallUpdate  bool
	}{
		// Success cases
		{
			name:              "Increase quantity success",
			userID:            10,
			productID:         100,
			quantity:          5,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			expectError:       false,
			shouldCallProduct: true,
			shouldCallUpdate:  true,
		},
		{
			name:              "Decrease quantity success",
			userID:            10,
			productID:         100,
			quantity:          1,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			expectError:       false,
			shouldCallProduct: true,
			shouldCallUpdate:  true,
		},
		// Set quantity to 0 removes item
		{
			name:              "Set quantity to 0 removes item",
			userID:            10,
			productID:         100,
			quantity:          0,
			expectError:       false,
			shouldCallProduct: false, // should not check product when quantity is 0
			shouldCallUpdate:  true,
		},
		// Invalid inputs
		{
			name:              "Invalid userID negative",
			userID:            -1,
			productID:         100,
			quantity:          1,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallProduct: false,
			shouldCallUpdate:  false,
		},
		{
			name:              "Invalid productID zero",
			userID:            10,
			productID:         0,
			quantity:          1,
			expectError:       true,
			expectedErrType:   ErrInvalidCartInput,
			shouldCallProduct: false,
			shouldCallUpdate:  false,
		},
		// Product not found
		{
			name:              "Product not found",
			userID:            10,
			productID:         999,
			quantity:          1,
			productErr:        pgx.ErrNoRows,
			expectError:       true,
			expectedErrType:   ErrProductNotFound,
			shouldCallProduct: true,
			shouldCallUpdate:  false,
		},
		// Product repository error
		{
			name:              "Product repository error",
			userID:            10,
			productID:         100,
			quantity:          1,
			productErr:        errors.New("database error"),
			expectError:       true,
			shouldCallProduct: true,
			shouldCallUpdate:  false,
		},
		// Insufficient stock
		{
			name:              "Insufficient stock for new quantity",
			userID:            10,
			productID:         100,
			quantity:          20,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 5},
			expectError:       true,
			expectedErrType:   ErrInsufficientStock,
			shouldCallProduct: true,
			shouldCallUpdate:  false,
		},
		// Update repository error (item not found)
		{
			name:              "Item not found in cart",
			userID:            10,
			productID:         100,
			quantity:          2,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			updateErr:         pgx.ErrNoRows,
			expectError:       true,
			expectedErrType:   ErrCartItemNotFound,
			shouldCallProduct: true,
			shouldCallUpdate:  true,
		},
		// Update repository error (generic)
		{
			name:              "Update repository error",
			userID:            10,
			productID:         100,
			quantity:          2,
			product:           &domain.Product{Id: 100, Name: "Laptop", Price: 1500.0, Quantity: 10},
			updateErr:         errors.New("database error"),
			expectError:       true,
			shouldCallProduct: true,
			shouldCallUpdate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := &mockCartRepository{updateCartItemErr: tt.updateErr}
			mockProductRepo := &mockProductRepositoryForCart{
				getProductByIDResult: tt.product,
				getProductByIDErr:    tt.productErr,
			}

			service := NewCartService(mockCartRepo, mockProductRepo)
			err := service.UpdateCartItem(tt.userID, tt.productID, tt.quantity)

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

			if tt.shouldCallProduct && mockProductRepo.getProductByIDCalls == 0 {
				t.Fatalf("expected product repo to be called")
			}
			if !tt.shouldCallProduct && mockProductRepo.getProductByIDCalls != 0 {
				t.Fatalf("expected product repo not to be called, got %d calls", mockProductRepo.getProductByIDCalls)
			}

			if tt.shouldCallUpdate {
				if mockCartRepo.updateCartItemCalls.calls != 1 {
					t.Fatalf("expected UpdateCartItemQuantity called once, got %d", mockCartRepo.updateCartItemCalls.calls)
				}
			} else {
				if mockCartRepo.updateCartItemCalls.calls != 0 {
					t.Fatalf("expected UpdateCartItemQuantity not called, got %d calls", mockCartRepo.updateCartItemCalls.calls)
				}
			}
		})
	}
}

// ============================================================================
// TESTS: RemoveFromCart
// ============================================================================

// TestCartService_RemoveFromCart table-driven tests for item removal.
func TestCartService_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		productID       int64
		removeErr       error
		expectError     bool
		expectedErrType error
	}{
		// Success case
		{
			name:      "Remove existing item success",
			userID:    10,
			productID: 100,
			expectError: false,
		},
		// Invalid inputs
		{
			name:            "Invalid userID negative",
			userID:          -1,
			productID:       100,
			expectError:     true,
			expectedErrType: ErrInvalidCartInput,
		},
		{
			name:            "Invalid productID zero",
			userID:          10,
			productID:       0,
			expectError:     true,
			expectedErrType: ErrInvalidCartInput,
		},
		// Item not found
		{
			name:            "Item not found in cart",
			userID:          10,
			productID:       999,
			removeErr:       pgx.ErrNoRows,
			expectError:     true,
			expectedErrType: ErrCartItemNotFound,
		},
		// Repository error
		{
			name:        "Repository error",
			userID:      10,
			productID:   100,
			removeErr:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := &mockCartRepository{removeCartItemErr: tt.removeErr}
			mockProductRepo := &mockProductRepositoryForCart{}

			service := NewCartService(mockCartRepo, mockProductRepo)
			err := service.RemoveFromCart(tt.userID, tt.productID)

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
// TESTS: ClearCart
// ============================================================================

// TestCartService_ClearCart tests cart clearing functionality.
func TestCartService_ClearCart(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		clearErr        error
		expectError     bool
		expectedErrType error
	}{
		// Success case
		{
			name:        "Clear cart success",
			userID:      10,
			expectError: false,
		},
		// Invalid inputs
		{
			name:            "Invalid userID negative",
			userID:          -1,
			expectError:     true,
			expectedErrType: ErrInvalidCartInput,
		},
		{
			name:            "Invalid userID zero",
			userID:          0,
			expectError:     true,
			expectedErrType: ErrInvalidCartInput,
		},
		// Repository error
		{
			name:        "Repository error",
			userID:      10,
			clearErr:    errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := &mockCartRepository{clearCartErr: tt.clearErr}
			mockProductRepo := &mockProductRepositoryForCart{}

			service := NewCartService(mockCartRepo, mockProductRepo)
			err := service.ClearCart(tt.userID)

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
