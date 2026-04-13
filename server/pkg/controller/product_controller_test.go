package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
)

type mockProductService struct {
	updateStockResult *domain.Product
	updateStockErr    error

	capturedProductID int64
	capturedQuantity  int
}

func (m *mockProductService) Add(productCreate *dto.ProductCreate) error {
	return nil
}

func (m *mockProductService) GetProducts(query dto.ProductListQuery) ([]*domain.Product, int64, int) {
	return []*domain.Product{}, 0, 0
}

func (m *mockProductService) GetAllProducts() []*domain.Product {
	return []*domain.Product{}
}

func (m *mockProductService) GetAllProductsByCategory(category string) []*domain.Product {
	return []*domain.Product{}
}

func (m *mockProductService) UpdateProductById(updatedProduct *dto.ProductCreate, productId int64) error {
	return nil
}

func (m *mockProductService) UpdateStockById(productId int64, quantity int) (*domain.Product, error) {
	m.capturedProductID = productId
	m.capturedQuantity = quantity
	return m.updateStockResult, m.updateStockErr
}

func (m *mockProductService) DeleteById(productId int64) error {
	return nil
}

func TestProductController_UpdateProductStockById_Success(t *testing.T) {
	e := echo.New()
	mockService := &mockProductService{
		updateStockResult: &domain.Product{
			Id:       7,
			Name:     "Keyboard",
			Price:    99.99,
			Quantity: 15,
			Category: "Electronics",
		},
	}
	controller := NewProductController(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/products/7/stock", bytes.NewBufferString(`{"quantity":15}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("7")

	err := controller.UpdateProductStockById(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if mockService.capturedProductID != 7 {
		t.Fatalf("expected product id 7, got %d", mockService.capturedProductID)
	}
	if mockService.capturedQuantity != 15 {
		t.Fatalf("expected quantity 15, got %d", mockService.capturedQuantity)
	}

	var responseBody map[string]interface{}
	if err = json.Unmarshal(rec.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if int(responseBody["quantity"].(float64)) != 15 {
		t.Fatalf("expected response quantity 15, got %v", responseBody["quantity"])
	}
}

func TestProductController_UpdateProductStockById_InvalidProductID(t *testing.T) {
	e := echo.New()
	controller := NewProductController(&mockProductService{})

	req := httptest.NewRequest(http.MethodPatch, "/products/abc/stock", bytes.NewBufferString(`{"quantity":5}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("abc")

	err := controller.UpdateProductStockById(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProductController_UpdateProductStockById_ProductNotFound(t *testing.T) {
	e := echo.New()
	mockService := &mockProductService{updateStockErr: service.ErrProductNotFound}
	controller := NewProductController(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/products/999/stock", bytes.NewBufferString(`{"quantity":1}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("999")

	err := controller.UpdateProductStockById(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestProductController_UpdateProductStockById_InternalError(t *testing.T) {
	e := echo.New()
	mockService := &mockProductService{updateStockErr: errors.New("database down")}
	controller := NewProductController(mockService)

	req := httptest.NewRequest(http.MethodPatch, "/products/7/stock", bytes.NewBufferString(`{"quantity":4}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("7")

	err := controller.UpdateProductStockById(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
