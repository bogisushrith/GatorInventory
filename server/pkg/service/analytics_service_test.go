package service

import (
	"testing"
	"time"

	"ims-intro/pkg/service/dto"
)

type mockAnalyticsRepository struct {
	dateFrom          *time.Time
	lowStockThreshold int
	topLimit          int
	lowStockLimit     int

	summaryResult  dto.AnalyticsSummary
	revenueResult  []dto.RevenueTrendPoint
	ordersResult   []dto.OrderTrendPoint
	topResult      []dto.TopProductMetric
	lowStockResult []dto.LowStockProduct
}

func (repositoryMock *mockAnalyticsRepository) EnsureSeedData() error {
	return nil
}

func (repositoryMock *mockAnalyticsRepository) GetSummary(dateFrom *time.Time, lowStockThreshold int) (dto.AnalyticsSummary, error) {
	repositoryMock.dateFrom = dateFrom
	repositoryMock.lowStockThreshold = lowStockThreshold
	return repositoryMock.summaryResult, nil
}

func (repositoryMock *mockAnalyticsRepository) GetRevenueTrend(dateFrom *time.Time) ([]dto.RevenueTrendPoint, error) {
	repositoryMock.dateFrom = dateFrom
	return repositoryMock.revenueResult, nil
}

func (repositoryMock *mockAnalyticsRepository) GetOrderTrend(dateFrom *time.Time) ([]dto.OrderTrendPoint, error) {
	repositoryMock.dateFrom = dateFrom
	return repositoryMock.ordersResult, nil
}

func (repositoryMock *mockAnalyticsRepository) GetTopProducts(dateFrom *time.Time, limit int) ([]dto.TopProductMetric, error) {
	repositoryMock.dateFrom = dateFrom
	repositoryMock.topLimit = limit
	return repositoryMock.topResult, nil
}

func (repositoryMock *mockAnalyticsRepository) GetLowStockProducts(threshold int, limit int) ([]dto.LowStockProduct, error) {
	repositoryMock.lowStockThreshold = threshold
	repositoryMock.lowStockLimit = limit
	return repositoryMock.lowStockResult, nil
}

func TestAnalyticsService_GetSummary_AppliesDefaults(t *testing.T) {
	repositoryMock := &mockAnalyticsRepository{summaryResult: dto.AnalyticsSummary{TotalOrders: 5}}
	serviceInstance := NewAnalyticsService(repositoryMock)

	summary, err := serviceInstance.GetSummary(dto.AnalyticsQuery{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if summary.TotalOrders != 5 {
		t.Fatalf("expected total orders 5, got %d", summary.TotalOrders)
	}
	if repositoryMock.lowStockThreshold != 10 {
		t.Fatalf("expected default low stock threshold 10, got %d", repositoryMock.lowStockThreshold)
	}
	if repositoryMock.dateFrom == nil {
		t.Fatal("expected default dateFrom for 30-day window")
	}
}

func TestAnalyticsService_GetTopProducts_UsesLimit(t *testing.T) {
	repositoryMock := &mockAnalyticsRepository{topResult: []dto.TopProductMetric{{ProductID: 1}}}
	serviceInstance := NewAnalyticsService(repositoryMock)

	result, err := serviceInstance.GetTopProducts(dto.AnalyticsQuery{Days: 7, TopProductsLimit: 3})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected one top product, got %d", len(result))
	}
	if repositoryMock.topLimit != 3 {
		t.Fatalf("expected top products limit 3, got %d", repositoryMock.topLimit)
	}
	if repositoryMock.dateFrom == nil {
		t.Fatal("expected dateFrom to be set")
	}
}

func TestAnalyticsService_GetOrderTrend_InvalidDays(t *testing.T) {
	repositoryMock := &mockAnalyticsRepository{}
	serviceInstance := NewAnalyticsService(repositoryMock)

	_, err := serviceInstance.GetOrderTrend(dto.AnalyticsQuery{Days: 999})
	if err == nil {
		t.Fatal("expected validation error for invalid days")
	}
	if err != ErrInvalidAnalyticsInput {
		t.Fatalf("expected ErrInvalidAnalyticsInput, got %v", err)
	}
}
