package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
)

type UserAnalyticsController struct {
	userAnalyticsService service.IUserAnalyticsService
}

func NewUserAnalyticsController(userAnalyticsService service.IUserAnalyticsService) *UserAnalyticsController {
	return &UserAnalyticsController{userAnalyticsService: userAnalyticsService}
}

func (controller *UserAnalyticsController) RegisterUserAnalyticsRoutes(e *echo.Echo) {
	group := e.Group("/user/analytics")
	group.Use(middleware.AuthMiddleware)
	group.Use(middleware.Authorize([]string{"user"}))

	group.GET("/summary", controller.GetSummary)
	group.GET("/recent-orders", controller.GetRecentOrders)
	group.GET("/top-products", controller.GetTopProducts)
	group.GET("/spending-trend", controller.GetSpendingTrend)
	group.GET("/recommendations", controller.GetRecommendations)
}

func (controller *UserAnalyticsController) GetSummary(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	summary, err := controller.userAnalyticsService.GetSummary(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch user analytics summary"))
	}

	return c.JSON(http.StatusOK, response.ToUserAnalyticsSummaryResponse(summary))
}

func (controller *UserAnalyticsController) GetRecentOrders(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	items, err := controller.userAnalyticsService.GetRecentOrders(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch recent orders"))
	}

	return c.JSON(http.StatusOK, response.ToUserRecentOrderResponse(items))
}

func (controller *UserAnalyticsController) GetTopProducts(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	items, err := controller.userAnalyticsService.GetTopProducts(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch top products"))
	}

	return c.JSON(http.StatusOK, response.ToUserTopProductResponse(items))
}

func (controller *UserAnalyticsController) GetSpendingTrend(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	items, err := controller.userAnalyticsService.GetSpendingTrend(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch spending trend"))
	}

	return c.JSON(http.StatusOK, response.ToRevenuePointResponse(items))
}

func (controller *UserAnalyticsController) GetRecommendations(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	items, err := controller.userAnalyticsService.GetRecommendations(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch recommendations"))
	}

	return c.JSON(http.StatusOK, response.ToUserRecommendationResponse(items))
}
