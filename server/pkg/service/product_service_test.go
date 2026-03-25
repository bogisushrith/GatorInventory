package service

import (
	"errors"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
	"testing"
)

/*
UNIT TESTING CONCEPTS (BACKEND - SERVICE LAYER)

1) What is unit testing in backend?
   - Unit testing verifies one small unit of code in isolation (here: ProductService methods).
   - We test business rules without running real database or HTTP server.

2) Why test service layer instead of controller or DB first?
   - Service layer contains business logic (validation, pagination math, decision rules).
   - Controller tests are more about HTTP request/response wiring.
   - DB tests are integration concerns and are slower/flakier for quick feedback.

3) What is mocking and why do we use it?
   - A mock is a fake implementation of a dependency (repository) with controlled outputs.
   - We use mocks to make tests deterministic, fast, and independent from real DB state.

4) Unit vs Integration vs End-to-End (E2E)
   - Unit test: one function/class in isolation with mocks.
   - Integration test: multiple layers work together (often with real DB).
   - E2E test: full app flow (UI/API/DB) from user perspective.
*/

type mockProductRepository struct {
	addErr error

	capturedAddedProduct *domain.Product
	addCalled            bool

	getProductsResult []*domain.Product
	getProductsTotal  int64
	getProductsErr    error
	capturedQuery     dto.ProductListQuery

	checkExistenceErr error
}

func (mockRepo *mockProductRepository) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, error) {
	mockRepo.capturedQuery = query
	return mockRepo.getProductsResult, mockRepo.getProductsTotal, mockRepo.getProductsErr
}

func (mockRepo *mockProductRepository) GetAllProducts() []*domain.Product {
	return []*domain.Product{}
}

func (mockRepo *mockProductRepository) GetProductsByCategory(category string) []*domain.Product {
	return []*domain.Product{}
}

func (mockRepo *mockProductRepository) AddProduct(product *domain.Product) error {
	mockRepo.addCalled = true
	mockRepo.capturedAddedProduct = product
	return mockRepo.addErr
}

func (mockRepo *mockProductRepository) CheckProductExistence(productId int64) error {
	return mockRepo.checkExistenceErr
}

func (mockRepo *mockProductRepository) UpdateProductById(updatedProduct *domain.Product, productId int64) error {
	return nil
}

func (mockRepo *mockProductRepository) DeleteProductById(productId int64) error {
	return nil
}

// TestProductService_Add_Success tests the happy path for creating a product.
// Why important: validates core business flow from validated DTO -> domain model -> repository call.
// Validates: service/business logic and proper mapping of input fields.
// Setup: build valid input and mock repository that succeeds.
// Execution: call service.Add with valid product data.
// Assertion: ensure no error, repo method called, and mapped product fields are correct.
func TestProductService_Add_Success(t *testing.T) {
	mockRepo := &mockProductRepository{}
	service := NewProductService(mockRepo)

	productCreate := &dto.ProductCreate{
		Name:     "Laptop",
		Price:    1500.50,
		Quantity: 3,
		Category: "Electronics",
	}

	err := service.Add(productCreate)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !mockRepo.addCalled {
		t.Fatalf("expected AddProduct to be called")
	}
	if mockRepo.capturedAddedProduct == nil {
		t.Fatalf("expected capturedAddedProduct to be set")
	}
	if mockRepo.capturedAddedProduct.Name != "Laptop" {
		t.Fatalf("expected name Laptop, got %s", mockRepo.capturedAddedProduct.Name)
	}
	if mockRepo.capturedAddedProduct.Category != "Electronics" {
		t.Fatalf("expected category Electronics, got %s", mockRepo.capturedAddedProduct.Category)
	}
}

// TestProductService_Add_EdgeCase_InvalidInput checks validation for empty/invalid payload.
// Why important: prevents bad data from entering database by enforcing service rules.
// Validates: validation business logic before repository interaction.
// Setup: input has empty name (invalid edge case).
// Execution: call service.Add.
// Assertion: expect validation error and ensure repository is NOT called.
func TestProductService_Add_EdgeCase_InvalidInput(t *testing.T) {
	mockRepo := &mockProductRepository{}
	service := NewProductService(mockRepo)

	productCreate := &dto.ProductCreate{
		Name:     "",
		Price:    99.99,
		Quantity: 2,
		Category: "Accessories",
	}

	err := service.Add(productCreate)

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
	if mockRepo.addCalled {
		t.Fatalf("expected AddProduct not to be called when validation fails")
	}
}

// TestProductService_Add_Failure_RepositoryError checks error propagation from repository.
// Why important: ensures service does not hide data layer failures.
// Validates: failure handling in service layer.
// Setup: valid input, mock repository configured to return an error.
// Execution: call service.Add.
// Assertion: expect same failure signal from service.
func TestProductService_Add_Failure_RepositoryError(t *testing.T) {
	mockRepo := &mockProductRepository{addErr: errors.New("database insert failed")}
	service := NewProductService(mockRepo)

	productCreate := &dto.ProductCreate{
		Name:     "Monitor",
		Price:    300,
		Quantity: 10,
		Category: "Electronics",
	}

	err := service.Add(productCreate)

	if err == nil {
		t.Fatalf("expected repository error, got nil")
	}
}

// TestProductService_GetProducts_Success_WithPagination tests pagination math + query pass-through.
// Why important: confirms service computes total pages correctly and keeps filter/pagination query intact.
// Validates: business logic (totalPages calculation) and orchestration logic.
// Setup: mock repository returns total=23 with limit=10 and sample products.
// Execution: call service.GetProducts with filters and pagination.
// Assertion: verify products, total count, total pages (ceil(23/10)=3), and forwarded query values.
func TestProductService_GetProducts_Success_WithPagination(t *testing.T) {
	mockRepo := &mockProductRepository{
		getProductsResult: []*domain.Product{
			{Id: 1, Name: "Laptop", Category: "Electronics"},
			{Id: 2, Name: "Mouse", Category: "Electronics"},
		},
		getProductsTotal: 23,
	}
	service := NewProductService(mockRepo)

	query := dto.ProductListQuery{
		Page:     1,
		Limit:    10,
		Search:   "lap",
		Category: "Electronics",
	}

	products, total, totalPages := service.GetProducts(query)

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if total != 23 {
		t.Fatalf("expected total 23, got %d", total)
	}
	if totalPages != 3 {
		t.Fatalf("expected totalPages 3, got %d", totalPages)
	}
	if mockRepo.capturedQuery.Search != "lap" {
		t.Fatalf("expected query search lap, got %s", mockRepo.capturedQuery.Search)
	}
	if mockRepo.capturedQuery.Category != "Electronics" {
		t.Fatalf("expected category Electronics, got %s", mockRepo.capturedQuery.Category)
	}
}

// TestProductService_GetProducts_EdgeCase_EmptyData checks behavior when repository returns no rows.
// Why important: empty states are common in dashboards and should not break API behavior.
// Validates: service handling of empty datasets.
// Setup: mock repository returns empty list and total=0.
// Execution: call service.GetProducts.
// Assertion: expect empty result, zero total, and zero total pages.
func TestProductService_GetProducts_EdgeCase_EmptyData(t *testing.T) {
	mockRepo := &mockProductRepository{
		getProductsResult: []*domain.Product{},
		getProductsTotal:  0,
	}
	service := NewProductService(mockRepo)

	query := dto.ProductListQuery{Page: 1, Limit: 10}

	products, total, totalPages := service.GetProducts(query)

	if len(products) != 0 {
		t.Fatalf("expected empty products list, got %d items", len(products))
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
	if totalPages != 0 {
		t.Fatalf("expected totalPages 0, got %d", totalPages)
	}
}

// TestProductService_GetProducts_Failure_RepositoryError verifies fallback output on repo failure.
// Why important: keeps service response predictable even when data layer fails.
// Validates: failure branch logic in GetProducts.
// Setup: mock repository returns an error.
// Execution: call service.GetProducts.
// Assertion: expect empty list and zero metadata as defined by current service implementation.
func TestProductService_GetProducts_Failure_RepositoryError(t *testing.T) {
	mockRepo := &mockProductRepository{getProductsErr: errors.New("database query failed")}
	service := NewProductService(mockRepo)

	query := dto.ProductListQuery{Page: 1, Limit: 10}

	products, total, totalPages := service.GetProducts(query)

	if len(products) != 0 {
		t.Fatalf("expected empty products list on failure, got %d items", len(products))
	}
	if total != 0 {
		t.Fatalf("expected total 0 on failure, got %d", total)
	}
	if totalPages != 0 {
		t.Fatalf("expected totalPages 0 on failure, got %d", totalPages)
	}
}
