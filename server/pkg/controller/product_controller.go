package controller

import (
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/request"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
	"net/http"
	"strconv"
)

type ProductController struct {
	productService service.IProductService
}

func NewProductController(productService service.IProductService) *ProductController {
	return &ProductController{productService}
}

func (controller *ProductController) RegisterProductRoutes(e *echo.Echo) {
	productsGroup := e.Group("/products")
	productsGroup.Use(middleware.AuthMiddleware)

	productsGroup.GET("", controller.GetAllProducts, middleware.Authorize([]string{"admin", "user"}))
	productsGroup.POST("", controller.AddNewProduct, middleware.Authorize([]string{"admin"}))
	productsGroup.PUT("/:id", controller.UpdateProductById, middleware.Authorize([]string{"admin"}))
	productsGroup.DELETE("/:id", controller.DeleteProductById, middleware.Authorize([]string{"admin"}))
}

func (controller *ProductController) GetAllProducts(c echo.Context) error {
	page := 1
	limit := 5

	if rawPage := c.QueryParam("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if rawLimit := c.QueryParam("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if limit > 50 {
		limit = 50
	}

	query := dto.ProductListQuery{
		Page:     page,
		Limit:    limit,
		Search:   c.QueryParam("search"),
		Category: c.QueryParam("category"),
	}

	if rawMinPrice := c.QueryParam("min_price"); rawMinPrice != "" {
		parsedMinPrice, err := strconv.ParseFloat(rawMinPrice, 64)
		if err == nil {
			query.MinPrice = &parsedMinPrice
		}
	}

	if rawMaxPrice := c.QueryParam("max_price"); rawMaxPrice != "" {
		parsedMaxPrice, err := strconv.ParseFloat(rawMaxPrice, 64)
		if err == nil {
			query.MaxPrice = &parsedMaxPrice
		}
	}

	products, total, totalPages := controller.productService.GetProducts(query)
	return c.JSON(http.StatusOK, response.ToPaginatedProductResponse(products, query.Page, query.Limit, total, totalPages))
}

func (controller *ProductController) AddNewProduct(c echo.Context) error {
	addProductResponse := new(request.AddProductRequest)

	err := c.Bind(addProductResponse)
	if err != nil || addProductResponse == nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the product structure"))
	}

	err = controller.productService.Add(addProductResponse.ToModel())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *ProductController) UpdateProductById(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no product id specified"))
	}

	productId, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: product id must be an integer"))
	}

	addProductResponse := new(request.AddProductRequest)
	err = c.Bind(addProductResponse)
	if err != nil || addProductResponse == nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the product structure"))
	}

	err = controller.productService.UpdateProductById(addProductResponse.ToModel(), int64(productId))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (controller *ProductController) DeleteProductById(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no product id specified"))
	}

	productId, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: product id must be an integer"))
	}

	err = controller.productService.DeleteById(int64(productId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
