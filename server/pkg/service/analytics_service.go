package service

import (
	"errors"
	"time"

	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

var ErrInvalidAnalyticsInput = errors.New("invalid analytics input")

type IAnalyticsService interface {
	GetSummary(query dto.AnalyticsQuery) (dto.AnalyticsSummary, error)
	GetRevenueTrend(query dto.AnalyticsQuery) ([]dto.RevenueTrendPoint, error)
	GetOrderTrend(query dto.AnalyticsQuery) ([]dto.OrderTrendPoint, error)
	GetTopProducts(query dto.AnalyticsQuery) ([]dto.TopProductMetric, error)
	GetLowStockProducts(query dto.AnalyticsQuery) ([]dto.LowStockProduct, error)
}

type AnalyticsService struct {
	analyticsRepository repository.IAnalyticsRepository
}

func NewAnalyticsService(analyticsRepository repository.IAnalyticsRepository) IAnalyticsService {
	return &AnalyticsService{analyticsRepository: analyticsRepository}
}

func (service *AnalyticsService) GetSummary(query dto.AnalyticsQuery) (dto.AnalyticsSummary, error) {
	normalizedQuery, err := normalizeAnalyticsQuery(query)
	if err != nil {
		return dto.AnalyticsSummary{}, err
	}
	if err = service.ensureSeedData(); err != nil {
		return dto.AnalyticsSummary{}, err
	}
	return service.analyticsRepository.GetSummary(dateFromFromDays(normalizedQuery.Days), normalizedQuery.LowStockThreshold)
}

func (service *AnalyticsService) GetRevenueTrend(query dto.AnalyticsQuery) ([]dto.RevenueTrendPoint, error) {
	normalizedQuery, err := normalizeAnalyticsQuery(query)
	if err != nil {
		return nil, err
	}
	if err = service.ensureSeedData(); err != nil {
		return nil, err
	}
	return service.analyticsRepository.GetRevenueTrend(dateFromFromDays(normalizedQuery.Days))
}

func (service *AnalyticsService) GetOrderTrend(query dto.AnalyticsQuery) ([]dto.OrderTrendPoint, error) {
	normalizedQuery, err := normalizeAnalyticsQuery(query)
	if err != nil {
		return nil, err
	}
	if err = service.ensureSeedData(); err != nil {
		return nil, err
	}
	return service.analyticsRepository.GetOrderTrend(dateFromFromDays(normalizedQuery.Days))
}

func (service *AnalyticsService) GetTopProducts(query dto.AnalyticsQuery) ([]dto.TopProductMetric, error) {
	normalizedQuery, err := normalizeAnalyticsQuery(query)
	if err != nil {
		return nil, err
	}
	if err = service.ensureSeedData(); err != nil {
		return nil, err
	}
	return service.analyticsRepository.GetTopProducts(dateFromFromDays(normalizedQuery.Days), normalizedQuery.TopProductsLimit)
}

func (service *AnalyticsService) GetLowStockProducts(query dto.AnalyticsQuery) ([]dto.LowStockProduct, error) {
	normalizedQuery, err := normalizeAnalyticsQuery(query)
	if err != nil {
		return nil, err
	}
	if err = service.ensureSeedData(); err != nil {
		return nil, err
	}
	return service.analyticsRepository.GetLowStockProducts(normalizedQuery.LowStockThreshold, normalizedQuery.LowStockItemsLimit)
}

func (service *AnalyticsService) ensureSeedData() error {
	return service.analyticsRepository.EnsureSeedData()
}

func normalizeAnalyticsQuery(query dto.AnalyticsQuery) (dto.AnalyticsQuery, error) {
	normalized := query

	if normalized.Days == 0 {
		normalized.Days = 30
	}
	if normalized.Days < 0 || normalized.Days > 365 {
		return dto.AnalyticsQuery{}, ErrInvalidAnalyticsInput
	}

	if normalized.LowStockThreshold == 0 {
		normalized.LowStockThreshold = 10
	}
	if normalized.LowStockThreshold < 1 || normalized.LowStockThreshold > 5000 {
		return dto.AnalyticsQuery{}, ErrInvalidAnalyticsInput
	}

	if normalized.TopProductsLimit == 0 {
		normalized.TopProductsLimit = 5
	}
	if normalized.TopProductsLimit < 1 || normalized.TopProductsLimit > 50 {
		return dto.AnalyticsQuery{}, ErrInvalidAnalyticsInput
	}

	if normalized.LowStockItemsLimit == 0 {
		normalized.LowStockItemsLimit = 10
	}
	if normalized.LowStockItemsLimit < 1 || normalized.LowStockItemsLimit > 200 {
		return dto.AnalyticsQuery{}, ErrInvalidAnalyticsInput
	}

	return normalized, nil
}

func dateFromFromDays(days int) *time.Time {
	if days <= 0 {
		return nil
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	return &cutoff
}
