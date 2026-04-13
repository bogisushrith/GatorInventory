package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/request"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
)

type CartController struct {
	cartService service.ICartService
}

func NewCartController(cartService service.ICartService) *CartController {
	return &CartController{cartService: cartService}
}

func (controller *CartController) RegisterCartRoutes(e *echo.Echo) {
	cartGroup := e.Group("/cart")
	cartGroup.Use(middleware.AuthMiddleware)
	cartGroup.Use(middleware.Authorize([]string{"user"}))

	cartGroup.GET("", controller.GetCart)
	cartGroup.POST("/add", controller.AddToCart)
	cartGroup.PATCH("/update", controller.UpdateCart)
	cartGroup.DELETE("/remove", controller.RemoveFromCart)
}

func (controller *CartController) GetCart(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	items, err := controller.cartService.GetCart(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCartInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ToCartItemResponseList(items))
}

func (controller *CartController) AddToCart(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	addRequest := new(request.AddCartItemRequest)
	if err := c.Bind(addRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind cart payload"))
	}

	err := controller.cartService.AddToCart(userID, int64(addRequest.ProductID), addRequest.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCartInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return c.JSON(http.StatusConflict, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *CartController) UpdateCart(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	updateRequest := new(request.UpdateCartItemRequest)
	if err := c.Bind(updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind cart payload"))
	}

	err := controller.cartService.UpdateCartItem(userID, int64(updateRequest.ProductID), updateRequest.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCartInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return c.JSON(http.StatusConflict, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrCartItemNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (controller *CartController) RemoveFromCart(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	removeRequest := new(request.RemoveCartItemRequest)
	if err := c.Bind(removeRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind cart payload"))
	}

	err := controller.cartService.RemoveFromCart(userID, int64(removeRequest.ProductID))
	if err != nil {
		if errors.Is(err, service.ErrInvalidCartInput) {
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
		}
		if errors.Is(err, service.ErrCartItemNotFound) {
			return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
