package service

import (
	"errors"

	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

var ErrInvalidUserAnalyticsInput = errors.New("invalid user analytics input")

type IUserAnalyticsService interface {
	GetSummary(userID int64) (dto.UserAnalyticsSummary, error)
	GetRecentOrders(userID int64) ([]dto.UserRecentOrder, error)
	GetTopProducts(userID int64) ([]dto.UserTopProduct, error)
	GetSpendingTrend(userID int64) ([]dto.RevenueTrendPoint, error)
	GetRecommendations(userID int64) ([]dto.UserRecommendation, error)
}

type UserAnalyticsService struct {
	userAnalyticsRepository repository.IUserAnalyticsRepository
}

func NewUserAnalyticsService(userAnalyticsRepository repository.IUserAnalyticsRepository) IUserAnalyticsService {
	return &UserAnalyticsService{userAnalyticsRepository: userAnalyticsRepository}
}

func (service *UserAnalyticsService) GetSummary(userID int64) (dto.UserAnalyticsSummary, error) {
	if userID <= 0 {
		return dto.UserAnalyticsSummary{}, ErrInvalidUserAnalyticsInput
	}
	return service.userAnalyticsRepository.GetSummary(userID)
}

func (service *UserAnalyticsService) GetRecentOrders(userID int64) ([]dto.UserRecentOrder, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserAnalyticsInput
	}
	return service.userAnalyticsRepository.GetRecentOrders(userID, 5)
}

func (service *UserAnalyticsService) GetTopProducts(userID int64) ([]dto.UserTopProduct, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserAnalyticsInput
	}
	return service.userAnalyticsRepository.GetTopProducts(userID, 5)
}

func (service *UserAnalyticsService) GetSpendingTrend(userID int64) ([]dto.RevenueTrendPoint, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserAnalyticsInput
	}
	return service.userAnalyticsRepository.GetSpendingTrend(userID)
}

func (service *UserAnalyticsService) GetRecommendations(userID int64) ([]dto.UserRecommendation, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserAnalyticsInput
	}
	return service.userAnalyticsRepository.GetRecommendations(userID, 6)
}
