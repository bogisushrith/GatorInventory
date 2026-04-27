package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
)

type analyticsControllerServiceMock struct {
	query dto.AnalyticsQuery

	summaryResult  dto.AnalyticsSummary
	summaryErr     error
	revenueResult  []dto.RevenueTrendPoint
	revenueErr     error
	ordersResult   []dto.OrderTrendPoint
	ordersErr      error
	topResult      []dto.TopProductMetric
	topErr         error
	lowStockResult []dto.LowStockProduct
	lowStockErr    error
}

func (m *analyticsControllerServiceMock) GetSummary(query dto.AnalyticsQuery) (dto.AnalyticsSummary, error) {
	m.query = query
	return m.summaryResult, m.summaryErr
}
func (m *analyticsControllerServiceMock) GetRevenueTrend(query dto.AnalyticsQuery) ([]dto.RevenueTrendPoint, error) {
	m.query = query
	return m.revenueResult, m.revenueErr
}
func (m *analyticsControllerServiceMock) GetOrderTrend(query dto.AnalyticsQuery) ([]dto.OrderTrendPoint, error) {
	m.query = query
	return m.ordersResult, m.ordersErr
}
func (m *analyticsControllerServiceMock) GetTopProducts(query dto.AnalyticsQuery) ([]dto.TopProductMetric, error) {
	m.query = query
	return m.topResult, m.topErr
}
func (m *analyticsControllerServiceMock) GetLowStockProducts(query dto.AnalyticsQuery) ([]dto.LowStockProduct, error) {
	m.query = query
	return m.lowStockResult, m.lowStockErr
}

func TestAnalyticsController_GetAnalyticsSummary_Success_ShouldReturnCounts(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{summaryResult: dto.AnalyticsSummary{TotalRevenue: 300, TotalOrders: 4, TotalProducts: 12, LowStockCount: 2}}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/summary?days=7&threshold=8", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetAnalyticsSummary(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if mockService.query.Days != 7 || mockService.query.LowStockThreshold != 8 {
		t.Fatalf("expected query values to be parsed, got %+v", mockService.query)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected json body, got %v", err)
	}
	if int(body["totalOrders"].(float64)) != 4 || int(body["lowStockCount"].(float64)) != 2 {
		t.Fatalf("unexpected summary payload: %+v", body)
	}
}

func TestAnalyticsController_GetRevenueTrend_Success_ShouldReturnSeries(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{revenueResult: []dto.RevenueTrendPoint{{Date: "2026-04-25", Value: 120}}}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/revenue-trend?days=30", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetRevenueTrend(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAnalyticsController_GetOrderTrend_Success_ShouldReturnSeries(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{ordersResult: []dto.OrderTrendPoint{{Date: "2026-04-25", Value: 5}}}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/order-trend?days=30", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetOrderTrend(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAnalyticsController_GetTopProducts_Success_ShouldReturnTable(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{topResult: []dto.TopProductMetric{{ProductID: 1, Name: "Laptop", TotalSold: 2, Revenue: 2400}}}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/top-products?days=30&top_limit=5", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetTopProducts(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAnalyticsController_GetLowStockProducts_Success_ShouldReturnTable(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{lowStockResult: []dto.LowStockProduct{{ID: 2, Name: "Mouse", Quantity: 3, Category: "Accessories"}}}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/low-stock?threshold=5&low_stock_limit=10", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetLowStockProducts(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAnalyticsController_InvalidQuery_ShouldFail(t *testing.T) {
	e := echo.New()
	controller := NewAnalyticsController(&analyticsControllerServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/analytics/summary?days=oops", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetAnalyticsSummary(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAnalyticsController_ServiceError_ShouldReturnServerError(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{summaryErr: errors.New("db down")}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/summary", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetAnalyticsSummary(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestAnalyticsController_InvalidAnalyticsInput_ShouldReturn400(t *testing.T) {
	e := echo.New()
	mockService := &analyticsControllerServiceMock{lowStockErr: service.ErrInvalidAnalyticsInput}
	controller := NewAnalyticsController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/analytics/low-stock?threshold=bad", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetLowStockProducts(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
