package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
)

type AnalyticsController struct {
	analyticsService service.IAnalyticsService
}

func NewAnalyticsController(analyticsService service.IAnalyticsService) *AnalyticsController {
	return &AnalyticsController{analyticsService: analyticsService}
}

func (controller *AnalyticsController) RegisterAnalyticsRoutes(e *echo.Echo) {
	admin := e.Group("/analytics")
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.Authorize([]string{"admin"}))

	admin.GET("/summary", controller.GetAnalyticsSummary)
	admin.GET("/revenue-trend", controller.GetRevenueTrend)
	admin.GET("/order-trend", controller.GetOrderTrend)
	admin.GET("/top-products", controller.GetTopProducts)
	admin.GET("/low-stock", controller.GetLowStockProducts)
}

func (controller *AnalyticsController) GetAnalyticsSummary(c echo.Context) error {
	query, err := buildAnalyticsQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	summary, err := controller.analyticsService.GetSummary(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid analytics query"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch analytics summary"))
	}

	return c.JSON(http.StatusOK, response.ToSummaryResponse(summary))
}

func (controller *AnalyticsController) GetRevenueTrend(c echo.Context) error {
	query, err := buildAnalyticsQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	trend, err := controller.analyticsService.GetRevenueTrend(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid analytics query"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch revenue trend"))
	}

	return c.JSON(http.StatusOK, response.ToRevenuePointResponse(trend))
}

func (controller *AnalyticsController) GetOrderTrend(c echo.Context) error {
	query, err := buildAnalyticsQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	trend, err := controller.analyticsService.GetOrderTrend(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid analytics query"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch order trend"))
	}

	return c.JSON(http.StatusOK, response.ToOrderPointResponse(trend))
}

func (controller *AnalyticsController) GetTopProducts(c echo.Context) error {
	query, err := buildAnalyticsQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	products, err := controller.analyticsService.GetTopProducts(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid analytics query"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch top products"))
	}

	return c.JSON(http.StatusOK, response.ToTopProductAnalyticsResponse(products))
}

func (controller *AnalyticsController) GetLowStockProducts(c echo.Context) error {
	query, err := buildAnalyticsQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	products, err := controller.analyticsService.GetLowStockProducts(query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAnalyticsInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid analytics query"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("failed to fetch low stock products"))
	}

	return c.JSON(http.StatusOK, response.ToLowStockProductResponse(products))
}

func buildAnalyticsQuery(c echo.Context) (dto.AnalyticsQuery, error) {
	query := dto.AnalyticsQuery{}

	if c.QueryParam("days") != "" {
		days, err := parseAnalyticsIntQuery(c, "days")
		if err != nil {
			return dto.AnalyticsQuery{}, err
		}
		query.Days = days
	}

	if c.QueryParam("threshold") != "" {
		threshold, err := parseAnalyticsIntQuery(c, "threshold")
		if err != nil {
			return dto.AnalyticsQuery{}, err
		}
		query.LowStockThreshold = threshold
	}

	if c.QueryParam("top_limit") != "" {
		limit, err := parseAnalyticsIntQuery(c, "top_limit")
		if err != nil {
			return dto.AnalyticsQuery{}, err
		}
		query.TopProductsLimit = limit
	}

	if c.QueryParam("low_stock_limit") != "" {
		limit, err := parseAnalyticsIntQuery(c, "low_stock_limit")
		if err != nil {
			return dto.AnalyticsQuery{}, err
		}
		query.LowStockItemsLimit = limit
	}

	return query, nil
}

func parseAnalyticsIntQuery(c echo.Context, key string) (int, error) {
	value := c.QueryParam(key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(key + " must be an integer")
	}
	return parsed, nil
}
