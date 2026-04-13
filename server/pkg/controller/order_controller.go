package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/request"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
)

type OrderController struct {
	orderService service.IOrderService
}

func NewOrderController(orderService service.IOrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (controller *OrderController) RegisterOrderRoutes(e *echo.Echo) {
	ordersGroup := e.Group("/orders")
	ordersGroup.Use(middleware.AuthMiddleware)

	ordersGroup.POST("", controller.CreateOrder, middleware.Authorize([]string{"user"}))
	ordersGroup.GET("", controller.GetAllOrders, middleware.Authorize([]string{"admin", "user"}))
	ordersGroup.GET("/:id", controller.GetOrderByID, middleware.Authorize([]string{"admin", "user"}))
}

func (controller *OrderController) CreateOrder(c echo.Context) error {
	createOrderRequest := new(request.CreateOrderRequest)
	if err := c.Bind(createOrderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind order payload"))
	}
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}

	orderID, err := controller.orderService.CreateOrder(userID, createOrderRequest.ToModel())
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrderInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return c.JSON(http.StatusConflict, response.NewErrorResponse(err.Error()))
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
			return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("orders schema is missing; run migrations"))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.CreateOrderResponse{OrderID: orderID})
}

func (controller *OrderController) GetAllOrders(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}
	role, _ := c.Get("role").(string)
	query, err := buildOrderQueryFromContext(c, userID, role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	orders, err := controller.orderService.GetOrders(userID, role, query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrderInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToOrderResponseList(orders))
}

func (controller *OrderController) GetOrderByID(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no order id specified"))
	}

	orderID, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: order id must be an integer"))
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("invalid user context"))
	}
	role, _ := c.Get("role").(string)

	order, err := controller.orderService.GetOrderByRole(userID, role, orderID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrderInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrOrderNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToOrderResponse(order))
}

func buildOrderQueryFromContext(c echo.Context, userID int64, role string) (dto.OrderListQuery, error) {
	query := dto.OrderListQuery{UserID: userID, Role: role}
	query.Search = strings.TrimSpace(c.QueryParam("search"))
	query.Status = strings.TrimSpace(c.QueryParam("status"))

	if rawDateFrom := strings.TrimSpace(c.QueryParam("date_from")); rawDateFrom != "" {
		parsedDateFrom, err := parseDateQueryValue(rawDateFrom)
		if err != nil {
			return dto.OrderListQuery{}, err
		}
		query.DateFrom = &parsedDateFrom
	}

	if rawDateTo := strings.TrimSpace(c.QueryParam("date_to")); rawDateTo != "" {
		parsedDateTo, err := parseDateQueryValue(rawDateTo)
		if err != nil {
			return dto.OrderListQuery{}, err
		}
		parsedDateTo = parsedDateTo.Add(24 * time.Hour)
		query.DateTo = &parsedDateTo
	}

	return query, nil
}

func parseDateQueryValue(rawValue string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", rawValue)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return parsedDate, nil
}
