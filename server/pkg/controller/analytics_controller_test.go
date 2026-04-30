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

type mockAnalyticsService struct {
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

func (serviceMock *mockAnalyticsService) GetSummary(query dto.AnalyticsQuery) (dto.AnalyticsSummary, error) {
	serviceMock.query = query
	return serviceMock.summaryResult, serviceMock.summaryErr
}

func (serviceMock *mockAnalyticsService) GetRevenueTrend(query dto.AnalyticsQuery) ([]dto.RevenueTrendPoint, error) {
	serviceMock.query = query
	return serviceMock.revenueResult, serviceMock.revenueErr
}

func (serviceMock *mockAnalyticsService) GetOrderTrend(query dto.AnalyticsQuery) ([]dto.OrderTrendPoint, error) {
	serviceMock.query = query
	return serviceMock.ordersResult, serviceMock.ordersErr
}

func (serviceMock *mockAnalyticsService) GetTopProducts(query dto.AnalyticsQuery) ([]dto.TopProductMetric, error) {
	serviceMock.query = query
	return serviceMock.topResult, serviceMock.topErr
}

func (serviceMock *mockAnalyticsService) GetLowStockProducts(query dto.AnalyticsQuery) ([]dto.LowStockProduct, error) {
	serviceMock.query = query
	return serviceMock.lowStockResult, serviceMock.lowStockErr
}

func TestAnalyticsController_GetAnalyticsSummary_Success(t *testing.T) {
	e := echo.New()
	serviceMock := &mockAnalyticsService{
		summaryResult: dto.AnalyticsSummary{
			TotalRevenue:  1337.42,
			TotalOrders:   12,
			TotalProducts: 33,
			LowStockCount: 4,
		},
	}

	controller := NewAnalyticsController(serviceMock)
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
	if serviceMock.query.Days != 7 {
		t.Fatalf("expected days 7, got %d", serviceMock.query.Days)
	}
	if serviceMock.query.LowStockThreshold != 8 {
		t.Fatalf("expected threshold 8, got %d", serviceMock.query.LowStockThreshold)
	}

	var body map[string]interface{}
	if err = json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if int(body["totalOrders"].(float64)) != 12 {
		t.Fatalf("expected totalOrders 12, got %v", body["totalOrders"])
	}
}

func TestAnalyticsController_GetRevenueTrend_InvalidQuery(t *testing.T) {
	e := echo.New()
	controller := NewAnalyticsController(&mockAnalyticsService{})
	req := httptest.NewRequest(http.MethodGet, "/analytics/revenue-trend?days=oops", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetRevenueTrend(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAnalyticsController_GetTopProducts_ServiceError(t *testing.T) {
	e := echo.New()
	serviceMock := &mockAnalyticsService{topErr: errors.New("db down")}
	controller := NewAnalyticsController(serviceMock)
	req := httptest.NewRequest(http.MethodGet, "/analytics/top-products", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := controller.GetTopProducts(ctx)
	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestAnalyticsController_GetLowStockProducts_InvalidAnalyticsInput(t *testing.T) {
	e := echo.New()
	serviceMock := &mockAnalyticsService{lowStockErr: service.ErrInvalidAnalyticsInput}
	controller := NewAnalyticsController(serviceMock)
	req := httptest.NewRequest(http.MethodGet, "/analytics/low-stock", nil)
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
